#!/bin/bash
# Copyright 2020 DSR Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -euo pipefail
source integration_tests/cli/common.sh

# Preparation of Actors
vid=$RANDOM
vid2=$RANDOM
vendor_account=vendor_account_$vid
second_vendor_account=vendor_account_$vid2
vendor_admin_account=vendor_admin_account
echo "Create First Vendor Account - $vendor_account"
# ACTION: Create a vendor account with a random Vendor ID
create_new_vendor_account $vendor_account $vid
echo "Create Second Vendor Account - $second_vendor_account"
# ACTION: Create a second vendor account with another random Vendor ID
create_new_vendor_account $second_vendor_account $vid2
echo "Create A VendorAdmin Account - $vendor_admin_account"
# ACTION: Create a VendorAdmin account for administrative operations
create_new_account $vendor_admin_account "VendorAdmin"

test_divider

# Query non existent
echo "Query non existent vendorinfo"
# ACTION: Query VendorInfo for a non-existent VID and expect it to not be found
result=$(dcld query vendorinfo vendor --vid=$vid)
check_response "$result" "Not Found"
echo "$result"

test_divider

echo "Request all vendor info must be empty"
# ACTION: Query all vendors and expect an empty list initially
result=$(dcld query vendorinfo all-vendors)
check_response "$result" "\[\]"
echo "$result"

test_divider

# Create a vendor info record
echo "Create VendorInfo Record for VID: $vid"
companyLegalName="XYZ IOT Devices Inc"
vendorName="XYZ Devices"
schema_version_0=0
# ACTION: Add VendorInfo record for the first VID using its vendor account
result=$(echo "test1234" | dcld tx vendorinfo add-vendor --vid=$vid --companyLegalName="$companyLegalName" --vendorName="$vendorName" --schemaVersion=$schema_version_0 --from=$vendor_account --yes)
result=$(get_txn_result "$result")
check_response "$result" "\"code\": 0"
echo "$result"

test_divider

# Query vendor info record
echo "Verify if VendorInfo Record for VID: $vid is present or not"
# ACTION: Query the newly created VendorInfo and verify its fields
result=$(dcld query vendorinfo vendor --vid=$vid)
check_response "$result" "\"vendorID\": $vid"
check_response "$result" "\"companyLegalName\": \"$companyLegalName\""
check_response "$result" "\"vendorName\": \"$vendorName\""
check_response "$result" "\"schemaVersion\": $schema_version_0"
echo "$result"

test_divider

echo "Request all vendor info"
# ACTION: Query all vendors and verify the created record is present
result=$(dcld query vendorinfo all-vendors)
check_response "$result" "\"vendorID\": $vid"
check_response "$result" "\"companyLegalName\": \"$companyLegalName\""
check_response "$result" "\"vendorName\": \"$vendorName\""
echo "$result"

test_divider

# Update vendor info with empty optional fields
echo "Update vendor info record for VID: $vid (with required fields only)"
# ACTION: Update VendorInfo record with only required fields (omitting optional ones)
result=$(dcld tx vendorinfo update-vendor --vid=$vid --from=$vendor_account --yes)
result=$(get_txn_result "$result")
check_response "$result" "\"code\": 0"
echo "$result"

test_divider

# Query updated vendor info record
echo "Verify that omitted fields in update object are not set to empty string"
# ACTION: Verify fields after partial update to ensure omitted fields weren't wiped
result=$(dcld query vendorinfo vendor --vid=$vid)
check_response "$result" "\"vendorID\": $vid"
check_response "$result" "\"companyLegalName\": \"$companyLegalName\""
check_response "$result" "\"vendorName\": \"$vendorName\""
echo "$result"

test_divider

# Update vendor info record
echo "Update vendor info record for VID: $vid"
companyLegalName="ABC Subsidiary Corporation"
vendorLandingPageURL="https://www.w3.org/"
# ACTION: Perform a full update of the VendorInfo record including optional fields
result=$(echo "test1234" | dcld tx vendorinfo update-vendor --vid=$vid --companyLegalName="$companyLegalName" --vendorLandingPageURL=$vendorLandingPageURL --vendorName="$vendorName" --schemaVersion=$schema_version_0 --from=$vendor_account --yes)
result=$(get_txn_result "$result")
check_response "$result" "\"code\": 0"
echo "$result"

test_divider

# Query updated vendor info record
echo "Verify if VendorInfo Record for VID: $vid is updated or not"
# ACTION: Verify the fields reflect the full update
result=$(dcld query vendorinfo vendor --vid=$vid)
check_response "$result" "\"vendorID\": $vid"
check_response "$result" "\"companyLegalName\": \"$companyLegalName\""
check_response "$result" "\"vendorName\": \"$vendorName\""
check_response "$result" "\"vendorLandingPageURL\": \"$vendorLandingPageURL\""
check_response "$result" "\"schemaVersion\": $schema_version_0"
echo "$result"

test_divider

# Create a vendor info record from a vendor account belonging to another vendor_account
vid1=$RANDOM
# ACTION: Attempt to add VendorInfo for a different VID and expect an authorization error
result=$(echo "test1234" | dcld tx vendorinfo add-vendor --vid=$vid1 --companyLegalName="$companyLegalName" --vendorName="$vendorName" --from=$vendor_account --yes 2>&1) || true
result=$(get_txn_result "$result")
echo "$result"
check_response_and_report "$result" "transaction should be signed by a vendor account associated with the vendorID $vid1"

test_divider

# Update a vendor info record from a vendor account belonging to another vendor_account
# ACTION: Attempt to update VendorInfo using an account from another vendor and expect an error
result=$(echo "test1234" | dcld tx vendorinfo update-vendor --vid=$vid --companyLegalName="$companyLegalName" --vendorName="$vendorName" --from=$second_vendor_account --yes 2>&1) || true
result=$(get_txn_result "$result")
echo "$result"
check_response_and_report "$result" "transaction should be signed by a vendor account associated with the vendorID $vid"

test_divider

# Create a vendor info reacord from a vendor admin account
echo "Create a vendor info reacord from a vendor admin account"
vid=$RANDOM
# ACTION: Create a VendorInfo record using a VendorAdmin account (allowed)
result=$(echo "test1234" | dcld tx vendorinfo add-vendor --vid=$vid --companyLegalName="$companyLegalName" --vendorName="$vendorName" --from=$vendor_admin_account --yes 2>&1) || true
result=$(get_txn_result "$result")
echo "$result"
check_response "$result" "\"code\": 0"
echo "$result"

test_divider

# Update the vendor info record by a vendor admin account
echo "Update the vendor info record by a vendor admin account"
companyLegalName1="New Corp"
vendorName1="New Vendor Name"
# ACTION: Update the VendorInfo record using a VendorAdmin account
result=$(echo "test1234" | dcld tx vendorinfo update-vendor --vid=$vid --companyLegalName="$companyLegalName1" --vendorName="$vendorName1" --from=$vendor_admin_account --yes)
result=$(get_txn_result "$result")
check_response "$result" "\"code\": 0"
echo "$result"

test_divider
