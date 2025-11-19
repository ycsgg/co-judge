import random
import os
from typing import List, Optional, Set, Tuple

class MipsGenerator:
    def __init__(self, max_instr: int = 400, seed: Optional[int] = None):
        if seed is not None:
            random.seed(seed)
        else:
            random.seed(os.urandom(16))
        
        self.max_instr = max_instr
        self.budget = max_instr
        self.lines: List[str] = []
        self.func_labels: List[str] = []  # 存储已生成的函数名，供 jal 跳转
        
        # 寄存器定义
        self.REG_T = [f"$t{i}" for i in range(10)]
        self.REG_S = [f"$s{i}" for i in range(8)]
        # $ra 需要小心使用，它是 jal 的隐式目标
        self.REG_ALL = self.REG_T + self.REG_S 
        self.ZERO = "$zero"
        self.RA = "$ra"

    def _emit(self, text: str):
        """输出指令并扣除预算"""
        if self.budget > 0:
            self.lines.append(text)
            # 伪指令或标签不严格扣预算，但为了防止无限生成，这里统一扣除
            # 只有空行或注释不扣
            if not text.strip().startswith("#") and ":" not in text:
                self.budget -= 1

    def _r_reg(self, exclude: Set[str] = None) -> str:
        """从通用寄存器池中随机选择一个，避开 exclude 中的寄存器"""
        if exclude is None:
            exclude = set()
        # 过滤掉被保护的寄存器
        candidates = [r for r in self.REG_ALL if r not in exclude]
        if not candidates:
            # 如果过滤完没了（极少情况），就回退到只用 $t0
            return "$t0"
        return random.choice(candidates)

    def _imm16(self) -> int:
        return random.randint(0, 0xFFFF)

    def _gen_arith(self, exclude_regs: Set[str]):
        """生成算术/逻辑指令：add, sub, ori, lui, nop"""
        op_type = random.random()
        rd = self._r_reg(exclude_regs)
        
        if op_type < 0.4: # R-type: add, sub
            rs = self._r_reg() # 源寄存器可以是任意值，不需要 exclude
            rt = self._r_reg()
            op = "add" if random.random() < 0.5 else "sub"
            self._emit(f"{op} {rd}, {rs}, {rt}")
        elif op_type < 0.7: # I-type: ori
            rs = self._r_reg()
            imm = self._imm16()
            self._emit(f"ori {rd}, {rs}, 0x{imm:04X}")
        elif op_type < 0.85: # lui
            imm = self._imm16()
            self._emit(f"lui {rd}, 0x{imm:04X}")
        else:
            self._emit("nop")

    def _gen_mem(self, exclude_regs: Set[str], base: str = "$s0"):
        """生成访存指令：lw, sw"""
        rt = self._r_reg(exclude_regs) # lw 的目标或 sw 的源
        # 偏移量 0-1020, 4字节对齐
        offset = random.randrange(0, 1021, 4)
        op = "lw" if random.random() < 0.5 else "sw"
        self._emit(f"{op} {rt}, {offset}({base})")

    def _gen_jal(self):
        """生成函数调用 jal"""
        if self.func_labels:
            target = random.choice(self.func_labels)
            self._emit(f"jal {target}")
            # jal 会修改 $ra，如果上下文对此敏感需额外处理
            # 本生成器假设 main 和 loop 都在顶层，不依赖 $ra 存活

    def _mix_body(self, count: int, exclude_regs: Set[str], allow_mem=True, allow_jal=True):
        """混合生成各种指令的主体"""
        for _ in range(count):
            if self.budget <= 0: break
            
            rand_val = random.random()
            # 10% 概率生成 jal (如果有可用函数)
            if allow_jal and self.func_labels and rand_val < 0.1:
                self._gen_jal()
            # 30% 概率访存
            elif allow_mem and rand_val < 0.4:
                self._gen_mem(exclude_regs)
            # 60% 概率算术
            else:
                self._gen_arith(exclude_regs)


    def _mk_loop(self, name: str, iters: int, body_ops: int, nested: bool = False):
        """
        生成循环结构。
        关键修改：将循环计数器加入 exclude_regs，防止循环体修改它。
        """
        # 挑选两个专用的寄存器作为计数器和步长，确保不在 REG_S 中 (假设 s0 是基址)
        # 为了安全，我们固定使用高位 t 寄存器，或者动态选择
        ctr_reg = "$t8"
        step_reg = "$t9"
        
        # 必须保护这两个寄存器不被循环体内的 add/sub/lw 覆盖
        current_excludes = {ctr_reg, step_reg, self.ZERO}
        
        self._emit(f"# --- begin loop {name} ---")
        self._emit(f"ori {ctr_reg}, {self.ZERO}, 0x{iters & 0xFFFF:04X}")
        self._emit(f"ori {step_reg}, {self.ZERO}, 0x0001")
        
        self._emit(f"{name}:")
        self._emit(f"beq {ctr_reg}, {self.ZERO}, {name}_end")
        
        # 生成循环体，传入黑名单
        self._mix_body(body_ops, exclude_regs=current_excludes, allow_mem=True)
        
        # 循环变量更新
        self._emit(f"sub {ctr_reg}, {ctr_reg}, {step_reg}")
        self._emit(f"beq {self.ZERO}, {self.ZERO}, {name}")
        self._emit(f"{name}_end:")
        self._emit(f"# --- end loop {name} ---")

    def _gen_leaf_functions(self, count: int = 2):
        """生成叶子函数（不调用其他函数的函数），供 jal 使用"""
        # 先把主程序跳过这些函数定义区，或者把它们放在 text 段末尾
        # 这里我们先生成，稍后在 generate() 中拼接到最后
        func_lines = []
        # 临时切换 self.lines 指向 func_lines，为了复用 _emit
        original_lines = self.lines
        self.lines = func_lines

        for i in range(count):
            fname = f"func_{i}"
            self.func_labels.append(fname)
            self._emit(f"\n{fname}:  # Leaf function")
            # 函数体内不能改 $s0 (数据基址)，也不能改 $ra (返回地址)
            # 函数体内可以随意改 $t 寄存器
            excludes = {"$s0", "$ra", self.ZERO}
            
            # 函数体短一点
            self._mix_body(random.randint(3, 6), exclude_regs=excludes, allow_jal=False)
            
            self._emit(f"jr {self.RA}")
        
        # 恢复
        self.lines = original_lines
        return func_lines

    def generate(self, level: int = 2):
        # 1. 头部配置
        self._emit(".text")
        # 初始化基址 $s0
        self._emit(f"lui $s0, 0x0000")
        self._emit(f"ori $s0, $s0, 0x0000")
        
        # 2. 预先生成一些函数体（存在内存里，最后打印）
        #    这样主逻辑里就可以 generate jal func_x 了
        func_code_lines = self._gen_leaf_functions(count=2)

        # 3. 初始化一些寄存器
        for r in random.sample(self.REG_T[:8], k=4):
            self._emit(f"ori {r}, {self.ZERO}, 0x{self._imm16():04X}")

        # 4. 线性混合代码
        self._mix_body(10, exclude_regs={"$s0", self.ZERO})

        # 5. 循环生成
        # Level 决定循环的个数和嵌套
        loop_count = 1 if level == 1 else 2
        if level >= 3: loop_count = 3
        
        for i in range(loop_count):
            if self.budget <= 0: break
            iters = random.randint(4, 10)
            ops = random.randint(5, 10)
            self._mk_loop(f"Loop_{i}", iters, ops)

        # 6. 结尾
        # 防止主程序如果不小心滑落到函数定义区
        # 通常模拟器运行到最后一行指令后应该停止，或者死循环
        self._emit(f"beq {self.ZERO}, {self.ZERO}, END_MAIN") 
        
        # 7. 附加函数定义代码到末尾
        self.lines.append("\n# --- Subroutines ---")
        self.lines.extend(func_code_lines)
        
        self._emit("\nEND_MAIN:")

        return self.lines

if __name__ == "__main__":
    # 使用示例
    generator = MipsGenerator(max_instr=300)
    # 参数 level 控制复杂度
    code = generator.generate(level=2)
    
    print("\n".join(code))
