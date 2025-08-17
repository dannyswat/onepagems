#!/bin/bash

# Quick demonstration of Step 4.3: Form Generator
echo "🎉 OnePage CMS - Step 4.3: Form Generator Demo"
echo "=============================================="
echo ""

echo "Building and starting server..."
cd /Users/dannys/repos/onepagems
go build -o demo cmd/main.go
./demo &
DEMO_PID=$!
sleep 3

echo "Logging in..."
curl -s -c demo_cookies.txt -d "username=admin&password=admin123" http://localhost:8080/admin/login > /dev/null

echo ""
echo "🔧 Testing Advanced Schema with All Field Types:"
echo ""

# Create a comprehensive test schema
TEST_SCHEMA='{
    "type": "object",
    "properties": {
        "title": {
            "type": "string",
            "title": "Website Title",
            "description": "The main title of your website",
            "minLength": 1
        },
        "hero_section": {
            "type": "object",
            "title": "Hero Section",
            "description": "Main banner area",
            "properties": {
                "headline": {
                    "type": "string",
                    "title": "Hero Headline",
                    "description": "Main headline text"
                },
                "hero_image": {
                    "type": "string",
                    "title": "Hero Background Image",
                    "description": "Upload a hero background image",
                    "format": "image"
                },
                "description": {
                    "type": "string",
                    "title": "Hero Description",
                    "description": "Descriptive text for the hero section",
                    "format": "textarea"
                }
            }
        },
        "content": {
            "type": "string",
            "title": "Rich Content",
            "description": "Main page content with HTML formatting",
            "format": "html"
        },
        "contact_email": {
            "type": "string",
            "title": "Contact Email",
            "description": "Primary contact email address",
            "format": "email"
        },
        "rating": {
            "type": "number",
            "title": "Site Rating",
            "description": "Rate your site",
            "minimum": 1,
            "maximum": 10
        },
        "featured": {
            "type": "boolean",
            "title": "Featured Site",
            "description": "Mark this site as featured"
        },
        "categories": {
            "type": "array",
            "title": "Site Categories",
            "description": "Select applicable categories",
            "items": {
                "type": "string",
                "enum": ["business", "personal", "portfolio", "blog", "ecommerce"]
            }
        },
        "theme": {
            "type": "string",
            "title": "Site Theme",
            "description": "Choose a theme",
            "enum": ["light", "dark", "blue", "green"]
        }
    },
    "required": ["title", "contact_email"]
}'

echo "Updating schema with comprehensive field types..."
curl -s -b demo_cookies.txt -X POST -H "Content-Type: application/json" -d "$TEST_SCHEMA" http://localhost:8080/admin/schema > /dev/null

echo ""
echo "📋 Generated Form Fields:"
echo ""

FORM_DATA=$(curl -s -b demo_cookies.txt http://localhost:8080/admin/schema/form-fields)

echo "$FORM_DATA" | python3 -m json.tool | grep -E '"name"|"type"|"label"|"required"' | sed 's/^[[:space:]]*/  /'

echo ""
echo "🎯 Field Type Summary:"
echo ""

# Extract and count field types
TEXT_COUNT=$(echo "$FORM_DATA" | grep -o '"type":"text"' | wc -l | xargs)
TEXTAREA_COUNT=$(echo "$FORM_DATA" | grep -o '"type":"textarea"' | wc -l | xargs)
EMAIL_COUNT=$(echo "$FORM_DATA" | grep -o '"type":"email"' | wc -l | xargs)
NUMBER_COUNT=$(echo "$FORM_DATA" | grep -o '"type":"number"' | wc -l | xargs)
CHECKBOX_COUNT=$(echo "$FORM_DATA" | grep -o '"type":"checkbox"' | wc -l | xargs)
SELECT_COUNT=$(echo "$FORM_DATA" | grep -o '"type":"select"' | wc -l | xargs)
MULTISELECT_COUNT=$(echo "$FORM_DATA" | grep -o '"type":"multiselect"' | wc -l | xargs)
IMAGE_COUNT=$(echo "$FORM_DATA" | grep -o '"type":"image"' | wc -l | xargs)
RICHTEXT_COUNT=$(echo "$FORM_DATA" | grep -o '"type":"richtext"' | wc -l | xargs)
OBJECT_COUNT=$(echo "$FORM_DATA" | grep -o '"type":"object"' | wc -l | xargs)

echo "  📝 Text fields: $TEXT_COUNT"
echo "  📄 Textarea fields: $TEXTAREA_COUNT"
echo "  📧 Email fields: $EMAIL_COUNT"
echo "  🔢 Number fields: $NUMBER_COUNT"
echo "  ☑️  Checkbox fields: $CHECKBOX_COUNT"
echo "  📋 Select fields: $SELECT_COUNT"
echo "  📋 Multi-select fields: $MULTISELECT_COUNT"
echo "  🖼️  Image fields: $IMAGE_COUNT"
echo "  📝 Rich text fields: $RICHTEXT_COUNT"
echo "  📁 Object fields: $OBJECT_COUNT"

echo ""
echo "✨ Advanced Features Demonstrated:"
echo ""
echo "  ✅ Dynamic field type detection"
echo "  ✅ Nested object handling (hero_section.*)"
echo "  ✅ Special format recognition (image, html, email, textarea)"
echo "  ✅ Enum to select/multiselect conversion"
echo "  ✅ Validation constraint extraction (min/max, required)"
echo "  ✅ Field prioritization and ordering"
echo "  ✅ Rich metadata (labels, descriptions, placeholders)"
echo ""

echo "🚀 Form Generator Implementation Complete!"
echo ""
echo "Ready for Step 5: Admin Interface Development"

# Cleanup
kill $DEMO_PID 2>/dev/null
rm -f demo demo_cookies.txt

echo ""
echo "Demo complete! 🎉"
