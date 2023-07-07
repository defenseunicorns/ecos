resource "aws_instance" "ecos_k8s" {
  instance_type = var.k8s_type
  ami           = var.amis[var.aws_region]

  key_name = aws_key_pair.generated_k8s_key.key_name

  subnet_id              = aws_subnet.private_subnet.0.id
  vpc_security_group_ids = [aws_security_group.ecos_k8s_sg.id]

  root_block_device {
    encrypted   = true
    volume_type = "gp2"
    volume_size = 100
  }

  tags = {
    Name = "ecos-k8s"
  }
}

output "k8s_private_ip" {
  value = aws_instance.ecos_k8s.private_ip
}
