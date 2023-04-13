terraform {
backend "local" {
  path = "../dev.tfstate"
}
#   cloud {
#     organization = "group24ece404"

#     workspaces {
#       name = "group24dev"
#     }
#   }
}

provider "google" {
  #credentials = "Users/tonyni/.config/gcloud/application_default_credentials.json"
  project     = "group24ece404"
  region      = "us-east1"
}

# Create new storage bucket in the US multi-region
# with coldline storage and settings for main_page_suffix and not_found_page
resource "random_id" "bucket_prefix" {
  byte_length = 8
}

resource "google_storage_bucket" "static_website" {
  name          = "ece461-dev.tonyni.ca"
  location      = "US"
  #storage_class = "COLDLINE"
  #uniform_bucket_level_access = true
  website {
    main_page_suffix = "index.html"
  }
}
# Make bucket public
resource "google_storage_bucket_iam_member" "member" {
  provider = google
  bucket   = google_storage_bucket.static_website.name
  role     = "roles/storage.objectViewer"
  member   = "allUsers"
}

# reserve IP address
resource "google_compute_global_address" "dev_cdn_ip" {
  name = "dev-cdn"
}

# backend bucket with CDN policy with default ttl settings
resource "google_compute_backend_bucket" "default" {
  name        = "ece461-comput-bucket"
  description = "Contains project frontend"
  bucket_name = google_storage_bucket.static_website.name
  enable_cdn  = true
  cdn_policy {
    cache_mode        = "CACHE_ALL_STATIC"
    client_ttl        = 3600
    default_ttl       = 3600
    max_ttl           = 86400
    negative_caching  = true
    serve_while_stale = 86400
  }
}

# url map
resource "google_compute_url_map" "default" {
  name            = "http-lb"
  default_service = google_compute_backend_bucket.default.id
}

# http proxy
resource "google_compute_target_http_proxy" "default" {
  name    = "http-lb-proxy"
  url_map = google_compute_url_map.default.id
}

resource "google_compute_global_forwarding_rule" "default" {
  name                  = "http-lb-forwarding-rule"
  ip_protocol           = "TCP"
  load_balancing_scheme = "EXTERNAL"
  port_range            = "80"
  target                = google_compute_target_http_proxy.default.id
  ip_address            = google_compute_global_address.dev_cdn_ip.id
}

# Upload a simple index.html page to the bucket
resource "google_storage_bucket_object" "indexHtml" {
  name         = "index.html"
  source      = "../../my-app/dist/index.html"
  content_type = "text/html"
  bucket       = google_storage_bucket.static_website.id
}

# Upload a simple index.html page to the bucket
resource "google_storage_bucket_object" "indexCSS" {
  name         = "assets/index-6e9558c7.css"
  source      = "../../my-app/dist/assets/index-6e9558c7.css"
  content_type = "text/css"
  bucket       = google_storage_bucket.static_website.id
}

# Upload a simple index.html page to the bucket
resource "google_storage_bucket_object" "indexJS" {
  name         = "assets/index-d8d261c8.js"
  source      = "../../my-app/dist/assets/index-d8d261c8.js"
  content_type = "text/javascript"
  bucket       = google_storage_bucket.static_website.id
}

output "cdn_ip_addr" {
  value = google_compute_global_address.dev_cdn_ip.address
}