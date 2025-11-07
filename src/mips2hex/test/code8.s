.text
addu $v1, $ra, $a0
L0:
sh $t6, -24167($s5)
sw $t5, -13726($a2)
jal END_LABEL
L1:
mflo $s5
jr $s3
jr $ra
mfhi $t4
L2:
jr $t1
mult $t9, $v0
L3:
divu $s5, $a2
jalr $s7, $t2
mflo $a0
mtlo $t5
L4:
mflo $v0
L5:
multu $a3, $s3
mflo $t9
divu $a3, $t8
mfhi $s7
mult $s6, $t6
mthi $t3
L6:
mflo $s4
mflo $t6
L7:
lh $ra, -21131($t5)
L8:
mflo $t5
divu $a1, $s5
mtlo $t4
mtlo $ra
divu $t0, $s1
mtlo $t0
mthi $t8
mflo $a0
jr $a1
mfhi $t2
mtlo $t7
mtlo $t4
mfhi $t4
divu $t8, $s6
blez $zero, L2
L9:
jr $a2
L10:
mthi $t6
multu $s6, $s3
divu $t1, $ra
L11:
mtlo $v0
multu $s0, $t0
div $s4, $s0
mflo $v0
divu $ra, $v1
mfhi $s5
multu $v1, $a1
srlv $s3, $a0, $s1
mfhi $s6
srl $s3, $s2, 20
mflo $ra
divu $a0, $s2
div $v1, $t4
divu $ra, $t0
jr $t5
mfhi $zero
mtlo $t9
div $zero, $t5
jr $s1
L12:
mflo $s2
mthi $t1
sltu $t8, $s2, $s1
multu $s3, $t1
addiu $t1, $t2, -5923
jr $t2
mtlo $t4
mflo $s1
mult $a0, $t2
div $s3, $t0
mult $s1, $t0
mthi $t9
divu $s3, $t7
jr $s5
addiu $t1, $s4, -29981
L13:
addiu $v1, $s5, 26062
mfhi $zero
and $t0, $t6, $a3
L14:
div $ra, $s0
mult $a3, $s7
L15:
divu $a2, $a0
mult $t9, $a3
div $s7, $ra
divu $s5, $t3
beq $a0, $ra, L0
L16:
multu $s4, $t5
mfhi $t4
mthi $a3
L17:
mtlo $t8
addi $a3, $s5, -24025
L18:
jr $s2
mthi $a2
sll $a2, $t7, 2
mtlo $a0
nop
jr $s3
multu $s4, $t9
divu $s5, $s4
sw $a1, -14349($s0)
divu $s4, $s7
divu $v0, $ra
jr $s3
mthi $t0
mtlo $s1
sltiu $t5, $t2, -12073
mthi $a3
L19:
and $t9, $s2, $t9
sllv $a3, $t9, $s5
jr $t2
sb $t3, -1785($t9)
div $t3, $s3
mtlo $t0
mflo $t4
mthi $s0
jal L18
mtlo $t8
mtlo $t8
L20:
mtlo $s0
mult $ra, $s0
L21:
li $t5, 0x8FECBEAD
mflo $s4
mthi $s5
multu $s6, $t5
jr $s1
mult $t5, $a0
srav $a3, $t5, $a0
L22:
div $a0, $s7
div $a1, $t1
xor $a3, $t3, $t6
multu $a0, $t1
divu $s4, $s5
mtlo $s3
divu $t1, $t7
L23:
mfhi $t1
mtlo $a1
L24:
mthi $s6
sb $a3, -5473($s2)
mthi $t6
divu $s7, $a3
srl $s3, $s6, 16
divu $t4, $a3
mflo $t0
mfhi $s1
mthi $ra
L25:
subu $v1, $t0, $s2
L26:
mflo $s5
L27:
multu $t4, $a1
mthi $ra
L28:
nop
L29:
nop
END_LABEL:
nop