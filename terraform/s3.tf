resource "aws_s3_bucket" "poker_audio_cache" {
  bucket = "poker-audio-cache-${var.region}"
}
