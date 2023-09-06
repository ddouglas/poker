resource "aws_acm_certificate" "poker" {
  domain_name       = local.default_domain
  validation_method = "DNS"

  subject_alternative_names = [
    "*.${local.default_domain}"
  ]

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_acm_certificate_validation" "poker" {
  certificate_arn         = aws_acm_certificate.poker.arn
  validation_record_fqdns = [for record in aws_route53_record.certificate_validation : record.fqdn]
}
