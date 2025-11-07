.text
beq $t4, $s1, END_LABEL
L0:
li $t6, 0xF43B49C7
L1:
multu $s4, $t3
L2:
jal END_LABEL
xor $s6, $s1, $a3
nor $a1, $t2, $s0
multu $t9, $t1
beq $t0, $s5, L1
jr $t5
lhu $s6, 85($zero)
bgtz $s3, END_LABEL
div $s4, $ra
bne $t6, $s0, L0
lbu $t9, -25472($t2)
slti $t0, $s2, -26908
L3:
beq $a1, $zero, L3
sll $a2, $a1, 13
addiu $v0, $ra, 11522
subu $t2, $t2, $s1
L4:
slti $a1, $t0, -30514
L5:
lbu $a2, -6368($zero)
jr $a1
lui $v0, 0x682F
mult $t6, $t8
slt $s4, $t0, $t9
j L3
j L4
beq $s0, $t6, L4
lb $t1, -29650($t2)
bgtz $a0, L4
li $a1, 0x23C22C0B
sh $t2, -18313($ra)
mflo $t5
sltiu $a0, $s6, -28148
lb $s2, -18231($t3)
div $t8, $a3
sltiu $t5, $s1, 22514
subu $s3, $t9, $v1
sllv $t0, $t9, $s4
L6:
subu $t9, $t0, $s0
multu $s4, $s3
mult $t6, $v0
L7:
sltu $a0, $ra, $s3
sll $s7, $t0, 4
li $s1, 0x4F0A7B1B
addu $s3, $s7, $s4
slti $zero, $v1, -24936
nop
mflo $t8
lw $t8, 13611($ra)
srl $t9, $zero, 31
bgtz $t3, L5
add $s3, $s4, $t3
addiu $t0, $s7, -20587
mfhi $s0
srav $t7, $s4, $a3
mfhi $t8
multu $s1, $ra
lb $t4, -25986($t9)
lhu $t8, -4832($t8)
xor $a0, $zero, $s6
L8:
lui $t2, 0x3FEF
L9:
mult $t6, $a2
L10:
mult $s5, $s7
L11:
addu $a1, $s2, $t7
lhu $t7, -8724($v1)
lhu $ra, 621($t0)
srav $a0, $t6, $a1
sltiu $zero, $t6, 26419
sh $s4, 17739($s2)
L12:
addu $t6, $ra, $t8
addi $s5, $s5, -1342
lb $t2, -8338($s7)
lb $ra, 26750($t7)
or $t7, $s2, $ra
L13:
and $t7, $s5, $a0
subu $s4, $t7, $t8
srl $a2, $t8, 17
mthi $s3
L14:
srl $s0, $s3, 26
blez $s4, L3
addi $v1, $t7, 3254
multu $t4, $t5
beq $t6, $t2, L8
add $ra, $v0, $t1
divu $t4, $a1
L15:
lbu $s6, -18864($t1)
sltiu $t8, $t4, -12700
xor $s1, $v0, $s3
beq $s2, $s7, L14
jr $v0
sb $s2, 30754($s0)
subu $t1, $ra, $s7
jalr $t1, $v1
srlv $t4, $s2, $zero
L16:
lbu $s5, -19416($a0)
jalr $t6, $s0
addiu $v1, $t8, 3342
lh $s2, -24116($t0)
jr $v0
xor $v1, $s3, $v1
sw $s0, 3675($s5)
lui $s5, 0x173E
lbu $t4, 13184($a3)
beq $t6, $a1, L6
multu $s1, $s0
lhu $s3, -30109($t2)
L17:
subu $a3, $t0, $a2
L18:
lw $t6, -23056($t6)
addiu $s0, $s7, 25701
bgtz $t5, L1
lw $s6, -22492($a0)
L19:
or $a3, $t3, $t1
addi $s6, $s3, 16284
L20:
lh $s6, 20560($t3)
multu $t4, $v1
divu $a0, $v1
mfhi $zero
bgtz $a1, L1
subu $s6, $a2, $s3
mtlo $t4
L21:
slt $a2, $a3, $s2
lbu $v0, 24132($ra)
j L20
jal L21
j L12
L22:
li $s0, 0xFB9FD6E7
sll $ra, $s6, 13
sll $s1, $a2, 11
L23:
bgtz $s7, L3
lw $v0, -16592($ra)
add $t2, $s3, $zero
lb $a0, 13645($t7)
nor $v0, $s3, $s4
nor $t4, $zero, $t5
L24:
srav $t6, $s2, $t2
and $t0, $a2, $s1
lhu $a2, -32061($s5)
divu $t0, $s7
L25:
bgtz $s3, L3
subu $s5, $s0, $t2
L26:
sllv $s0, $s4, $v0
slt $s7, $a3, $s4
L27:
srlv $s4, $t1, $s1
nop
L28:
mflo $t5
mtlo $t8
L29:
subu $v1, $v1, $a1
srav $s3, $t9, $s1
slt $s3, $t4, $v1
END_LABEL:
nop