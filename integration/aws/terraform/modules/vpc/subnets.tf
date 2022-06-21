locals {
  subnet_az_letters = ["a", "b", "c"]
  subnet_cidrs      = module.subnet_addrs.network_cidr_blocks
  subnet_azs        = [for s in local.subnet_az_letters : "${data.aws_region.current.name}${s}"]
  subnet_az_map     = { for s in local.subnet_az_letters : s => "${data.aws_region.current.name}${s}" }
}

resource "aws_subnet" "private" {
  for_each = local.subnet_az_map

  vpc_id                  = aws_vpc.vpc.id
  availability_zone       = each.value
  cidr_block              = local.subnet_cidrs["private-${each.key}"]
  map_public_ip_on_launch = false

  tags = {
    "test" : "vpc-${var.id}-subnet-${each.key}"
  }
}

resource "aws_route_table_association" "private" {
  for_each = local.subnet_az_map

  subnet_id      = aws_subnet.private[each.key].id
  route_table_id = aws_route_table.private.id
}

resource "aws_subnet" "public" {
  for_each = local.subnet_az_map

  vpc_id                  = aws_vpc.vpc.id
  availability_zone       = each.value
  cidr_block              = local.subnet_cidrs["public-${each.key}"]
  map_public_ip_on_launch = true

  tags = {
    "test" : "vpc-${var.id}-public-subnet-${each.key}"
  }
}

resource "aws_route_table_association" "public" {
  for_each = local.subnet_az_map

  subnet_id      = aws_subnet.public[each.key].id
  route_table_id = aws_route_table.public.id
}
