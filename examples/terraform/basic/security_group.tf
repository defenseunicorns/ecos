resource "aws_default_security_group" "default_sg" {
  vpc_id = aws_vpc.ecos_vpc.id

  tags = {
    Name = "ecos-default-sg"
  }
}

resource "aws_security_group" "ecos_bastion_sg" {
  name   = "ecos-bastion-sg"
  vpc_id = aws_vpc.ecos_vpc.id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "ecos-bastion-sg"
  }
}

resource "aws_security_group" "ecos_k8s_sg" {
  name   = "ecos-k8s-sg"
  vpc_id = aws_vpc.ecos_vpc.id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/24"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "ecos-k8s-sg"
  }
}
