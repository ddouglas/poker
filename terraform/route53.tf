# data "cloudflare_zone" "main" {
#   name = "onetwentyseven.dev"
# }

# resource "cloudflare_record" "nameserver" {
#   for_each = toset(aws_route53_zone.poker.name_servers)
#   zone_id  = data.cloudflare_zone.main.zone_id
#   name     = "poker"
#   value    = each.value
#   type     = "NS"
#   proxied  = false
# }


# resource "aws_route53_zone" "poker" {
#   name = local.default_domain
# }

# resource "aws_route53_record" "certificate_validation" {
#   for_each = {
#     for dvo in aws_acm_certificate.poker.domain_validation_options : dvo.domain_name => {
#       name   = dvo.resource_record_name
#       record = dvo.resource_record_value
#       type   = dvo.resource_record_type
#     }
#   }

#   allow_overwrite = true
#   name            = each.value.name
#   records         = [each.value.record]
#   ttl             = 60
#   type            = each.value.type
#   zone_id         = aws_route53_zone.poker.zone_id
# }

# output "zone_ns_records" {
#   value = aws_route53_zone.poker.name_servers
# }

# resource "aws_route53_record" "api" {
#   name    = aws_apigatewayv2_domain_name.poker.domain_name
#   type    = "A"
#   zone_id = aws_route53_zone.poker.zone_id

#   alias {
#     name                   = aws_apigatewayv2_domain_name.poker.domain_name_configuration[0].target_domain_name
#     zone_id                = aws_apigatewayv2_domain_name.poker.domain_name_configuration[0].hosted_zone_id
#     evaluate_target_health = false
#   }
# }
