# NAT gateway
resource "aws_eip" "nat_gw_eip" {
  domain = "vpc"
}

resource "aws_nat_gateway" "nat_gw" {
  allocation_id = aws_eip.nat_gw_eip.id
  subnet_id     = aws_subnet.public_subnet.0.id

  tags = {
    Name = "ecos-nat-gw"
  }
}

# Private subnet
resource "aws_subnet" "private_subnet" {
  count = length(var.azs)

  vpc_id     = aws_vpc.ecos_vpc.id
  cidr_block = "10.0.${length(var.azs) + count.index}.0/24"

  tags = {
    Name = "ecos-private-${var.azs[count.index]}"
  }
}

# Route table for the private subnet
resource "aws_route_table" "private_rte" {
  vpc_id = aws_vpc.ecos_vpc.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.nat_gw.id
  }

  tags = {
    Name = "ecos-private-rte"
  }
}

# Associate the private route table with the private subnets
resource "aws_route_table_association" "private_assoc" {
  count = length(var.azs)

  subnet_id      = aws_subnet.private_subnet[count.index].id
  route_table_id = aws_route_table.private_rte.id
}

