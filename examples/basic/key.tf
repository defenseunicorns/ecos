#Bastion Key
resource "tls_private_key" "bastion" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "generated_bastion_key" {
  key_name   = "ecos-bastion-key"
  public_key = tls_private_key.bastion.public_key_openssh
}

output "private_bastion_key" {
  value     = tls_private_key.bastion.private_key_pem
  sensitive = true
}

#K8s Key
resource "tls_private_key" "k8s" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "generated_k8s_key" {
  key_name   = "ecos-k8s-key"
  public_key = tls_private_key.k8s.public_key_openssh
}

output "private_k8s_key" {
  value     = tls_private_key.k8s.private_key_pem
  sensitive = true
}
