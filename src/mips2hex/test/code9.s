.text
bne $t7, $ra, END_LABEL
beq $a3, $t6, END_LABEL
slti $t3, $s5, -415
slt $a3, $a3, $s6
L0:
sb $t7, 15997($s6)
sw $s2, 239($ra)
bne $t9, $s7, L0
bne $s6, $t0, END_LABEL
sltu $s2, $a3, $ra
slt $v1, $t4, $s2
L1:
beq $v0, $t3, END_LABEL
sltiu $a0, $s7, -16656
L2:
sltiu $s4, $t1, 21006
L3:
bne $v1, $v0, L1
sltiu $t4, $s0, 31800
lw $s3, -14526($t2)
bne $t4, $t8, END_LABEL
slt $t2, $t3, $t4
L4:
jr $s3
slti $t5, $t6, 2106
bne $zero, $a1, L2
sltiu $s7, $s6, 5660
sltiu $t0, $t5, -6252
slti $a1, $t5, -15849
bne $t2, $ra, END_LABEL
bne $t1, $s3, END_LABEL
slt $s3, $v1, $t2
sltu $s7, $a3, $t1
beq $s5, $t4, L0
slti $t2, $s1, 26234
sltiu $s3, $t9, -16610
bne $a1, $t1, L1
sltiu $a0, $v0, -3277
L5:
bgtz $v1, L5
sltu $s2, $t6, $s5
sltiu $s6, $t3, -3359
L6:
slti $s6, $a3, -20987
sra $s2, $s5, 26
L7:
slt $t2, $a0, $t2
slti $t6, $v0, 11919
blez $v1, END_LABEL
sltu $s2, $a0, $v0
slt $s7, $a1, $a1
sltu $a2, $t3, $s6
sltu $a0, $t9, $v0
L8:
slti $t8, $s7, -7471
L9:
sltiu $a2, $t6, 14186
slti $t1, $t2, -31202
bne $s0, $s0, END_LABEL
slt $t2, $t1, $t8
slt $s4, $s1, $s6
bne $t9, $a2, L7
L10:
beq $zero, $t4, L6
sltiu $a3, $a2, 8912
L11:
sltiu $t7, $t5, -18442
sltu $a2, $v1, $t1
L12:
slti $t2, $s7, 19070
L13:
beq $t6, $s4, L5
sltiu $s7, $t3, -14404
slti $s2, $a3, -12487
L14:
sltu $t0, $s6, $v0
beq $v0, $s6, L13
blez $a2, L12
bne $t0, $t5, L14
L15:
slti $v0, $t7, -7263
slt $t9, $t6, $t1
bne $t2, $s0, L0
sltiu $s6, $s6, 11913
L16:
slt $s1, $s0, $v1
bne $s7, $a0, L3
L17:
bne $s4, $t5, END_LABEL
beq $s0, $t7, L17
slt $t8, $t0, $s7
beq $t8, $a2, L2
sh $t2, 1986($s4)
sltiu $s5, $t3, -13415
slti $ra, $s7, 22508
beq $ra, $v1, L10
L18:
slt $v1, $t7, $t6
sltu $s0, $t0, $t2
and $t4, $s1, $a3
L19:
divu $t7, $s3
L20:
slti $t7, $s6, -4360
slt $v1, $v1, $t6
bne $a2, $a3, L16
bne $zero, $a2, L13
sltiu $s5, $t5, 26351
bne $s5, $t5, L13
slti $s2, $t7, -28643
slt $t3, $a3, $s6
slt $v0, $s6, $v1
L21:
slti $ra, $zero, 25289
sltiu $s6, $s1, -8104
bne $s5, $s1, L9
add $v0, $t6, $t9
sltu $s3, $s2, $t1
slt $s0, $t3, $t5
slti $v1, $v1, -462
L22:
beq $s4, $t4, L12
beq $s1, $t7, L16
beq $s7, $v0, L1
sltiu $t6, $s1, -17786
sltiu $t7, $s3, -20528
slti $s4, $s5, -27270
sltu $s7, $a2, $t1
mtlo $s6
slti $t6, $t9, -8385
L23:
bne $s7, $t6, L17
L24:
sltu $t5, $s4, $a2
sltu $s2, $a3, $s2
L25:
slt $v1, $ra, $s6
L26:
bne $t9, $t7, L25
srl $t4, $s7, 1
L27:
beq $s3, $t7, L11
beq $t1, $t8, L11
beq $v1, $a0, L27
slti $s0, $v0, 16394
divu $t8, $a2
srlv $s0, $s5, $v1
L28:
addu $t1, $a2, $s6
sltu $t0, $s0, $t5
slti $t1, $t3, 7623
bne $s7, $t1, L24
L29:
bgtz $s7, L7
slti $t9, $t8, -4971
bne $t7, $s0, L11
bgtz $t3, L3
sltiu $t5, $a2, 18139
and $t7, $a1, $t4
slti $s4, $s3, -19959
bne $s0, $v0, L20
sltiu $s4, $t5, 2945
slt $s0, $s7, $s5
sltiu $t8, $t3, -1953
bne $t5, $t4, L9
slt $zero, $v0, $s2
beq $s7, $s6, L10
slt $t6, $a3, $t3
multu $t5, $a0
sltiu $t2, $s7, -6595
sltiu $a1, $t1, 10701
bne $s3, $t2, L28
sltu $t8, $t3, $s6
bne $t3, $t2, L25
lhu $ra, -11373($t8)
slti $v1, $v1, -18606
sltiu $s1, $t4, -31618
slt $s0, $zero, $s4
slti $s5, $s6, -13323
bne $s5, $s6, L21
END_LABEL:
nop