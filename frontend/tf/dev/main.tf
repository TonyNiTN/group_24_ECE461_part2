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
  name          = "${random_id.bucket_prefix.hex}-static-website-bucket"
  location      = "US"
  #storage_class = "COLDLINE"
#   uniform_bucket_level_access = true
  website {
    main_page_suffix = "index.html"
  }
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
  name         = "assets/index-8cb29f15.css"
  source      = "../../my-app/dist/assets/index-8cb29f15.css"
  bucket       = google_storage_bucket.static_website.id
}

# Upload a simple index.html page to the bucket
resource "google_storage_bucket_object" "indexJS" {
  name         = "assets/index-80bbade2.js"
  source      = "../../my-app/dist/assets/index-80bbade2.js"
  bucket       = google_storage_bucket.static_website.id
}

resource "google_cloud_run_service" "bucketfwdsvc" {
  name     = "bucketfwd"
  location = "us-east1"

  template {
    spec{
        containers {
            image = "gcr.io/group24ece404/bucketfwd:0.1.1"
        }
    } 
  }
  autogenerate_revision_name = true
}

# Make bucket public
# resource "google_storage_bucket_iam_member" "member" {
#   provider = google
#   bucket   = google_storage_bucket.static_website.name
#   role     = "roles/storage.objectViewer"
#   member   = "allUsers"
# }