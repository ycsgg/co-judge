.text   # 程序入口地址（0x3000）
lui $s0, 0x0000
ori $s0, $s0, 0x0000
ori $t7, $zero, 0x7E84
ori $t5, $zero, 0x826D
ori $t2, $zero, 0x6260
ori $t1, $zero, 0x2F8C
ori $t8, $zero, 0x7793
ori $t6, $zero, 0x049B
add $s7, $t6, $s4
add $s2, $t1, $t6
sub $t4, $t0, $t5
ori $t8, $t9, 0xA4FA
ori $s6, $s5, 0x4BD8
ori $t1, $s0, 0x4F97
ori $s0, $t1, 0x2414
nop
sub $s0, $s4, $t5
ori $s4, $zero, 0x6188
ori $s1, $t9, 0x6ED2
sub $s3, $t0, $t1
ori $s0, $s7, 0x3615
sub $s2, $t8, $s7
ori $t6, $zero, 0x10A1
sub $s1, $t0, $t2
ori $s3, $t8, 0x71D1
sub $s0, $t9, $t0
nop
add $t2, $s6, $t2
nop
lw $t9, 620($s0)
lw $s7, 340($s0)
lw $s1, 800($s0)
lw $s4, 1000($s0)
sw $s3, 72($s0)
lw $s1, 408($s0)
lw $t4, 468($s0)
sw $s7, 52($s0)
sw $s0, 272($s0)
sw $t9, 736($s0)
# --- begin loop Loop1, 10 iterations ---
ori $t8, $zero, 0x000A
ori $t9, $zero, 0x0001
Loop1:
beq $t8, $zero, Loop1_end
sub $s1, $t5, $s7
add $s2, $s0, $t0
sub $s1, $t4, $s6
sub $t0, $s6, $t5
sw $t4, 572($s0)
lw $t3, 728($s0)
sw $t0, 460($s0)
sw $t1, 260($s0)
sub $t8, $t8, $t9
beq $zero, $zero, Loop1
Loop1_end:
# --- end loop Loop1 ---
# --- begin loop Loop2, 6 iterations ---
ori $t8, $zero, 0x0006
ori $t9, $zero, 0x0001
Loop2:
beq $t8, $zero, Loop2_end
ori $s0, $t7, 0xB032
ori $t5, $t1, 0x00DD
ori $t0, $s0, 0x64B7
sub $s3, $t2, $s1
ori $s4, $t5, 0xF29B
ori $s0, $t2, 0x9895
sub $t8, $t8, $t9
beq $zero, $zero, Loop2
Loop2_end:
# --- end loop Loop2 ---
# --- begin loop Loop3, 2 iterations ---
ori $t8, $zero, 0x0002
ori $t9, $zero, 0x0001
Loop3:
beq $t8, $zero, Loop3_end
ori $s3, $s4, 0xEC2E
ori $t5, $s2, 0x1089
sw $s4, 208($s0)
lw $t7, 552($s0)
sw $t7, 560($s0)
nop
sub $t8, $t8, $t9
beq $zero, $zero, Loop3
Loop3_end:
# --- end loop Loop3 ---
END:
nop
