output "vpc_id" {
  value = aws_vpc.vpc.id
}

output "private_subnet_ids" {
  value = toset([for s in aws_subnet.private : s.id])
}

output "private_subnet_azs" {
  value = values(local.subnet_az_map)
}

output "private_subnet_az_map" {
  value = { for letter, az in local.subnet_az_map : az => aws_subnet.private[letter].id }
}

output "default_sg_id" {
  value = aws_default_security_group.default.id
}
