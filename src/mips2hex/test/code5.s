.text
bne $s2, $t2, END_LABEL
L0:
bne $s7, $s6, L0
L1:
beq $a1, $t2, L1
blez $t9, L0
blez $t0, END_LABEL
beq $t3, $s0, END_LABEL
bne $s6, $s1, END_LABEL
L2:
beq $t2, $a3, L2
blez $v0, L2
beq $t4, $v0, L1
sub $t5, $t6, $s2
bne $s6, $s2, END_LABEL
bne $t2, $t9, L1
beq $a2, $t2, L0
L3:
mfhi $v1
sub $s0, $t0, $t6
divu $t2, $ra
blez $v0, L2
bgtz $s1, END_LABEL
beq $t8, $t3, L3
blez $zero, END_LABEL
L4:
bgtz $t6, L1
L5:
addu $zero, $s6, $ra
bgtz $ra, L1
bne $t8, $s7, L2
beq $ra, $s5, L2
bne $zero, $t3, L3
beq $t3, $t6, L5
L6:
bne $t0, $v0, END_LABEL
L7:
beq $s4, $s0, L6
bgtz $s3, L0
L8:
blez $t3, L1
bgtz $ra, L7
blez $ra, L8
L9:
bgtz $t5, L5
L10:
blez $a1, L5
beq $t5, $t4, L2
lw $t8, 23028($s4)
bgtz $v1, L6
add $t9, $t0, $a1
L11:
beq $t8, $s6, L1
bne $zero, $t7, L8
multu $s1, $v0
beq $a0, $t1, L4
L12:
sltiu $t3, $t7, 15199
beq $t7, $s5, L4
bgtz $t9, L2
L13:
blez $s0, L11
bne $s2, $s7, L2
L14:
bgtz $zero, L3
lhu $a1, 8161($t0)
beq $a3, $v0, L6
bgtz $t1, L8
bgtz $t2, L14
L15:
blez $a0, L3
L16:
bgtz $v0, L2
blez $t8, L2
L17:
bne $s2, $a0, L11
blez $v0, L2
blez $t5, L8
bgtz $v0, L9
L18:
bgtz $s5, L14
blez $s2, L7
blez $s2, L7
bgtz $zero, L0
bgtz $t3, END_LABEL
bne $s2, $s2, L1
bgtz $t3, L18
L19:
bne $s0, $a0, L3
beq $s6, $s2, L5
L20:
blez $s4, END_LABEL
bne $s1, $s3, L4
bgtz $s7, L0
bgtz $t0, L14
bgtz $t8, L12
bgtz $t1, L11
blez $t1, L15
bgtz $t7, L0
addiu $ra, $zero, 27671
bgtz $s7, L19
lb $v1, 26759($a0)
mflo $a1
bgtz $v1, L7
divu $v0, $t9
j L10
div $s5, $a2
L21:
bne $s3, $t7, L0
blez $s1, END_LABEL
beq $a3, $s1, L14
L22:
bne $t7, $s6, L15
L23:
sltiu $t9, $a2, -24061
L24:
bne $s3, $t0, L6
beq $t7, $v0, L12
bne $a1, $s5, L10
L25:
bgtz $t5, L14
blez $zero, L25
bgtz $s7, L3
bne $t5, $a3, L17
bne $a2, $t4, L15
L26:
bne $s2, $v0, L18
bne $t4, $t8, L0
L27:
blez $t0, L27
L28:
sltu $t1, $a1, $v1
bne $t0, $t2, L23
blez $a3, L9
blez $zero, L7
bgtz $t6, L10
L29:
bne $ra, $s1, L2
blez $t7, L9
blez $t4, L16
bne $zero, $s5, L8
beq $v0, $s4, L27
bgtz $a1, L24
bgtz $s0, L9
addiu $v0, $t5, -28131
bgtz $s3, L26
bgtz $v0, L27
blez $t9, L10
beq $t5, $t6, L19
mthi $t1
mult $s2, $ra
bgtz $t2, END_LABEL
blez $t1, L24
mfhi $a1
bgtz $t8, L24
blez $t2, L25
bgtz $a0, L27
divu $s2, $s0
bgtz $v1, L29
blez $t1, L16
bne $s0, $a0, L17
blez $s7, L16
beq $s4, $a2, END_LABEL
bgtz $t6, L9
bne $t7, $t4, L1
beq $t7, $t4, L19
beq $a2, $t8, L21
beq $v1, $t7, L23
and $s6, $t5, $s0
blez $a1, L25
lui $ra, 0x915A
bne $a0, $a0, L4
bgtz $s6, L5
addi $v1, $a1, 27099
beq $zero, $a3, L9
sltu $s0, $t1, $v1
sllv $t0, $t0, $t3
bgtz $t2, L22
jal L27
beq $t6, $t0, L10
END_LABEL:
nop