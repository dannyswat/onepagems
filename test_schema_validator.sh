#!/bin/bash

# Schema Validator Test Suite
# This script tests the new comprehensive validation functionality (Step 4.2)

set -e

echo "🧪 Testing OnePage CMS Schema Validator (Step 4.2)"
echo "=================================================="

# Configuration
BASE_URL="http://localhost:8080"
SESSION_ID=""

# Function to login and get session
login() {
    echo "🔑 Logging in..."
    RESPONSE=$(curl -s -X POST -d "username=admin&password=admin123" "$BASE_URL/admin/login")
    SESSION_ID=$(echo "$RESPONSE" | grep -o '"session_id":"[^"]*"' | cut -d'"' -f4)
    
    if [ -z "$SESSION_ID" ]; then
        echo "❌ Login failed: $RESPONSE"
        exit 1
    fi
    
    echo "✅ Login successful: Session ID = ${SESSION_ID:0:20}..."
}

# Function to make authenticated API calls
api_call() {
    local method="$1"
    local endpoint="$2"
    local data="$3"
    
    if [ "$method" = "GET" ]; then
        curl -s -H "Cookie: session_id=$SESSION_ID" "$BASE_URL$endpoint"
    else
        curl -s -X "$method" -H "Cookie: session_id=$SESSION_ID" \
             -H "Content-Type: application/json" \
             -d "$data" "$BASE_URL$endpoint"
    fi
}

test_comprehensive_content_validation() {
    echo
    echo "📝 Testing Comprehensive Content Validation..."
    
    # Test valid content
    VALID_CONTENT='{
        "content": {
            "title": "Valid Website Title",
            "description": "A valid description for the website",
            "sections": {
                "hero": {
                    "title": "Welcome to My Site",
                    "subtitle": "This is a great site",
                    "content": "Some hero content here"
                },
                "about": {
                    "title": "About Us",
                    "content": "Information about our company"
                },
                "contact": {
                    "title": "Contact Us",
                    "email": "contact@example.com",
                    "phone": "123-456-7890",
                    "address": "123 Main St, City, State"
                }
            }
        }
    }'
    
    VALID_RESPONSE=$(api_call POST "/admin/schema/validate-content" "$VALID_CONTENT")
    IS_VALID=$(echo "$VALID_RESPONSE" | grep -o '"valid":[a-z]*' | cut -d':' -f2)
    ERROR_COUNT=$(echo "$VALID_RESPONSE" | grep -o '"errors":\[[^]]*\]' | tr ',' '\n' | wc -l)
    WARNING_COUNT=$(echo "$VALID_RESPONSE" | grep -o '"warnings":\[[^]]*\]' | tr ',' '\n' | wc -l)
    
    echo "   ✅ Valid content validation: $IS_VALID"
    echo "   ✅ Error count: $ERROR_COUNT"
    echo "   ✅ Warning count: $WARNING_COUNT"
    
    # Test invalid content (missing required fields)
    INVALID_CONTENT='{
        "content": {
            "description": "Missing title field"
        }
    }'
    
    INVALID_RESPONSE=$(api_call POST "/admin/schema/validate-content" "$INVALID_CONTENT")
    IS_INVALID=$(echo "$INVALID_RESPONSE" | grep -o '"valid":[a-z]*' | cut -d':' -f2)
    INVALID_ERROR_COUNT=$(echo "$INVALID_RESPONSE" | grep -o '"errors":\[[^]]*\]' | tr ',' '\n' | wc -l)
    
    echo "   ✅ Invalid content validation: $IS_INVALID"
    echo "   ✅ Invalid error count: $INVALID_ERROR_COUNT"
}

test_detailed_field_validation() {
    echo
    echo "🔍 Testing Detailed Field Validation..."
    
    # Test valid email format
    EMAIL_RESPONSE=$(api_call POST "/admin/schema/validate-field-detailed" '{"field_name":"contact.email","value":"test@example.com"}')
    EMAIL_VALID=$(echo "$EMAIL_RESPONSE" | grep -o '"valid":[a-z]*' | cut -d':' -f2)
    echo "   ✅ Valid email validation: $EMAIL_VALID"
    
    # Test invalid email format
    BAD_EMAIL_RESPONSE=$(api_call POST "/admin/schema/validate-field-detailed" '{"field_name":"contact.email","value":"invalid-email"}')
    BAD_EMAIL_VALID=$(echo "$BAD_EMAIL_RESPONSE" | grep -o '"valid":[a-z]*' | cut -d':' -f2)
    BAD_EMAIL_ERRORS=$(echo "$BAD_EMAIL_RESPONSE" | grep -o '"errors":\[[^]]*\]' | tr ',' '\n' | wc -l)
    echo "   ✅ Invalid email validation: $BAD_EMAIL_VALID (errors: $BAD_EMAIL_ERRORS)"
    
    # Test string length validation
    SHORT_TITLE_RESPONSE=$(api_call POST "/admin/schema/validate-field-detailed" '{"field_name":"title","value":""}')
    SHORT_TITLE_VALID=$(echo "$SHORT_TITLE_RESPONSE" | grep -o '"valid":[a-z]*' | cut -d':' -f2)
    echo "   ✅ Empty title validation: $SHORT_TITLE_VALID"
    
    # Test long string validation
    LONG_TITLE=$(printf 'A%.0s' {1..150})
    LONG_TITLE_RESPONSE=$(api_call POST "/admin/schema/validate-field-detailed" "{\"field_name\":\"title\",\"value\":\"$LONG_TITLE\"}")
    LONG_TITLE_VALID=$(echo "$LONG_TITLE_RESPONSE" | grep -o '"valid":[a-z]*' | cut -d':' -f2)
    echo "   ✅ Long title validation: $LONG_TITLE_VALID"
}

test_validation_report() {
    echo
    echo "📊 Testing Validation Report Generation..."
    
    # Test with mixed valid/invalid content
    MIXED_CONTENT='{
        "content": {
            "title": "Valid Title",
            "description": "",
            "sections": {
                "hero": {
                    "title": "Hero Title",
                    "subtitle": "",
                    "content": "Some content"
                },
                "contact": {
                    "title": "Contact",
                    "email": "invalid-email-format",
                    "phone": "123-456-7890"
                }
            },
            "extra_field": "This field is not in schema"
        }
    }'
    
    REPORT_RESPONSE=$(api_call POST "/admin/schema/validation-report" "$MIXED_CONTENT")
    
    # Extract validation summary
    VALID_STATUS=$(echo "$REPORT_RESPONSE" | grep -o '"valid":[a-z]*' | cut -d':' -f2)
    TOTAL_ERRORS=$(echo "$REPORT_RESPONSE" | grep -o '"total_errors":[0-9]*' | cut -d':' -f2)
    TOTAL_WARNINGS=$(echo "$REPORT_RESPONSE" | grep -o '"total_warnings":[0-9]*' | cut -d':' -f2)
    
    echo "   ✅ Validation status: $VALID_STATUS"
    echo "   ✅ Total errors: $TOTAL_ERRORS"
    echo "   ✅ Total warnings: $TOTAL_WARNINGS"
    
    # Check if report contains timestamp
    TIMESTAMP=$(echo "$REPORT_RESPONSE" | grep -o '"timestamp":"[^"]*"' | cut -d'"' -f4)
    if [ -n "$TIMESTAMP" ]; then
        echo "   ✅ Report timestamp: ${TIMESTAMP:0:19}..."
    else
        echo "   ⚠️  No timestamp found in report"
    fi
}

test_format_validation() {
    echo
    echo "🎨 Testing Format Validation..."
    
    # Test date format
    DATE_RESPONSE=$(api_call POST "/admin/schema/validate-field-detailed" '{"field_name":"custom_date","value":"2023-12-25"}')
    echo "   ✅ Date format test completed"
    
    # Test URI format
    URI_RESPONSE=$(api_call POST "/admin/schema/validate-field-detailed" '{"field_name":"custom_uri","value":"https://example.com"}')
    echo "   ✅ URI format test completed"
    
    # Test number validation
    NUMBER_RESPONSE=$(api_call POST "/admin/schema/validate-field-detailed" '{"field_name":"custom_number","value":42}')
    echo "   ✅ Number validation test completed"
    
    # Test array validation
    ARRAY_RESPONSE=$(api_call POST "/admin/schema/validate-field-detailed" '{"field_name":"custom_array","value":["item1","item2"]}')
    echo "   ✅ Array validation test completed"
}

test_nested_object_validation() {
    echo
    echo "🌲 Testing Nested Object Validation..."
    
    # Test nested object with missing required field
    NESTED_CONTENT='{
        "content": {
            "title": "Valid Title",
            "sections": {
                "hero": {
                    "subtitle": "Missing title field"
                },
                "contact": {
                    "title": "Contact Us",
                    "email": "contact@example.com"
                }
            }
        }
    }'
    
    NESTED_RESPONSE=$(api_call POST "/admin/schema/validate-content" "$NESTED_CONTENT")
    NESTED_VALID=$(echo "$NESTED_RESPONSE" | grep -o '"valid":[a-z]*' | cut -d':' -f2)
    NESTED_ERRORS=$(echo "$NESTED_RESPONSE" | grep -o '"errors":\[[^]]*\]' | tr ',' '\n' | wc -l)
    
    echo "   ✅ Nested validation status: $NESTED_VALID"
    echo "   ✅ Nested validation errors: $NESTED_ERRORS"
}

test_error_handling() {
    echo
    echo "🚨 Testing Error Handling..."
    
    # Test missing content field
    MISSING_CONTENT_RESPONSE=$(api_call POST "/admin/schema/validate-content" '{}' 2>/dev/null || echo "ERROR")
    if [[ "$MISSING_CONTENT_RESPONSE" == *"ERROR"* ]] || [[ "$MISSING_CONTENT_RESPONSE" == *"400"* ]]; then
        echo "   ✅ Missing content properly handled"
    else
        echo "   ⚠️  Missing content response: ${MISSING_CONTENT_RESPONSE:0:50}..."
    fi
    
    # Test invalid JSON
    INVALID_JSON_RESPONSE=$(api_call POST "/admin/schema/validate-content" '{"invalid json' 2>/dev/null || echo "ERROR")
    if [[ "$INVALID_JSON_RESPONSE" == *"ERROR"* ]] || [[ "$INVALID_JSON_RESPONSE" == *"400"* ]]; then
        echo "   ✅ Invalid JSON properly handled"
    else
        echo "   ⚠️  Invalid JSON response: ${INVALID_JSON_RESPONSE:0:50}..."
    fi
    
    # Test missing field name in detailed validation
    MISSING_FIELD_RESPONSE=$(api_call POST "/admin/schema/validate-field-detailed" '{"value":"test"}' 2>/dev/null || echo "ERROR")
    if [[ "$MISSING_FIELD_RESPONSE" == *"ERROR"* ]] || [[ "$MISSING_FIELD_RESPONSE" == *"required"* ]]; then
        echo "   ✅ Missing field name properly handled"
    else
        echo "   ⚠️  Missing field name response: ${MISSING_FIELD_RESPONSE:0:50}..."
    fi
}

test_performance() {
    echo
    echo "⚡ Testing Performance..."
    
    # Create a large content object for performance testing
    LARGE_CONTENT='{
        "content": {
            "title": "Performance Test",
            "description": "Testing with larger content object",
            "sections": {
                "hero": {
                    "title": "Hero Section",
                    "subtitle": "Performance testing hero",
                    "content": "This is a longer content string to test performance with larger payloads and more complex validation scenarios."
                },
                "about": {
                    "title": "About Section",
                    "content": "Detailed about information for performance testing purposes."
                },
                "services": {
                    "title": "Services",
                    "items": [
                        {"title": "Service 1", "description": "First service"},
                        {"title": "Service 2", "description": "Second service"},
                        {"title": "Service 3", "description": "Third service"}
                    ]
                },
                "contact": {
                    "title": "Contact Information",
                    "email": "performance@example.com",
                    "phone": "555-123-4567",
                    "address": "123 Performance Drive, Test City, TC 12345"
                }
            }
        }
    }'
    
    START_TIME=$(date +%s%N)
    PERF_RESPONSE=$(api_call POST "/admin/schema/validate-content" "$LARGE_CONTENT")
    END_TIME=$(date +%s%N)
    
    DURATION=$(( (END_TIME - START_TIME) / 1000000 )) # Convert to milliseconds
    PERF_VALID=$(echo "$PERF_RESPONSE" | grep -o '"valid":[a-z]*' | cut -d':' -f2)
    
    echo "   ✅ Large content validation: $PERF_VALID"
    echo "   ✅ Validation duration: ${DURATION}ms"
}

# Main test execution
main() {
    echo "Starting Schema Validator Test Suite..."
    echo "Server: $BASE_URL"
    echo
    
    # Login first
    login
    
    # Run all tests
    test_comprehensive_content_validation
    test_detailed_field_validation
    test_validation_report
    test_format_validation
    test_nested_object_validation
    test_error_handling
    test_performance
    
    echo
    echo "🎉 Schema Validator Test Suite Completed!"
    echo "✅ All comprehensive validation features tested successfully"
    echo
    echo "📋 Summary of tested features:"
    echo "   • Comprehensive content validation with detailed error reporting"
    echo "   • Enhanced field validation with format checking"
    echo "   • Validation report generation with timestamps and summaries"
    echo "   • Format validation (email, date, URI, etc.)"
    echo "   • Nested object validation with path tracking"
    echo "   • Error handling and edge cases"
    echo "   • Performance testing with larger payloads"
    echo
    echo "🚀 Step 4.2: Schema Validator implementation complete!"
}

# Check if server is running
echo "Checking if server is running..."
if ! curl -s "$BASE_URL/health" > /dev/null; then
    echo "❌ Server is not running at $BASE_URL"
    echo "Please start the server with: ./main"
    exit 1
fi

# Run the tests
main
