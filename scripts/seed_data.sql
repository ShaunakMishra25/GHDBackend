-- Seed categories & products (run this after migrations)
-- Idempotent: only inserts if a category/product does not already exist.

-- Clean up duplicate Rice legacy rows
DELETE FROM products
WHERE name_en = 'Rice'
  AND id <> (SELECT MIN(id) FROM products WHERE name_en = 'Rice');

-- ---------- CATEGORIES ----------
INSERT INTO categories (name_en, name_hi, image_url, sort_order, is_active)
SELECT d.name_en, d.name_hi, d.image_url, d.sort_order, true
FROM (VALUES
    ('Grocery',               'किराना सामान',        'https://images.unsplash.com/photo-1542838132-92c53300491e?w=400',  1),
    ('Fruits & Vegetables',   'फल और सब्जियाँ',       'https://images.unsplash.com/photo-1610832958506-aa56368176cf?w=400',  2),
    ('Dairy & Eggs',          'डेयरी और अंडे',        'https://images.unsplash.com/photo-1550583724-b2692b85b150?w=400',  3),
    ('Bakery & Snacks',       'बेकरी और नाश्ता',       'https://images.unsplash.com/photo-1509440159596-0249088772ff?w=400',  4),
    ('Beverages',             'पेय पदार्थ',            'https://images.unsplash.com/photo-1544145945-f90425340c7e?w=400',  5),
    ('Personal Care',         'व्यक्तिगत देखभाल',      'https://images.unsplash.com/photo-1556228720-195a672e8a03?w=400',  6),
    ('Household Essentials',  'घरेलू सामान',           'https://images.unsplash.com/photo-1584568694244-14fbdf83bd30?w=400',  7)
) AS d(name_en, name_hi, image_url, sort_order)
WHERE NOT EXISTS (SELECT 1 FROM categories c WHERE c.name_en = d.name_en);

-- ---------- PRODUCTS ----------
-- Fruits & Vegetables
INSERT INTO products (category_id, name_en, name_hi, description_en, description_hi, price, unit, stock_qty, image_url)
SELECT (SELECT id FROM categories WHERE name_en = 'Fruits & Vegetables'),
       d.name_en, d.name_hi, d.description_en, d.description_hi, d.price, d.unit, d.stock_qty, d.image_url
FROM (VALUES
    ('Tomatoes',  'टमाटर',   'Fresh ripe tomatoes',      'ताजे पके टमाटर',   40.00, 'kg',    200, 'https://images.unsplash.com/photo-1546094096-0df4bcaaa337?w=400'),
    ('Onions',    'प्याज',   'Premium onions',           'बेहतरीन प्याज',     35.00, 'kg',    300, 'https://images.unsplash.com/photo-1518977956812-cd3dbada9e31?w=400'),
    ('Potatoes',  'आलू',     'Farm fresh potatoes',      'ताजे आलू',           30.00, 'kg',    300, 'https://images.unsplash.com/photo-1518977676601-b53f82aba655?w=400'),
    ('Bananas',   'केला',    'Sweet bananas',            'मीठे केले',          60.00, 'dozen', 150, 'https://images.unsplash.com/photo-1571771894821-ce9b6c11b08e?w=400'),
    ('Apples',    'सेब',     'Crisp apples',             'कुरकुरे सेब',        120.00, 'kg',   100, 'https://images.unsplash.com/photo-1560806887-1e4cd0b6cbd6?w=400'),
    ('Green Chilli','हरी मिर्च','Fresh green chillies',  'ताजी हरी मिर्च',     45.00, 'kg',    80,  'https://images.unsplash.com/photo-1518998053901-5348d3961a04?w=400')
) AS d(name_en, name_hi, description_en, description_hi, price, unit, stock_qty, image_url)
WHERE NOT EXISTS (SELECT 1 FROM products p WHERE p.name_en = d.name_en AND p.category_id = (SELECT id FROM categories WHERE name_en = 'Fruits & Vegetables'));

-- Dairy & Eggs
INSERT INTO products (category_id, name_en, name_hi, description_en, description_hi, price, unit, stock_qty, image_url)
SELECT (SELECT id FROM categories WHERE name_en = 'Dairy & Eggs'),
       d.name_en, d.name_hi, d.description_en, d.description_hi, d.price, d.unit, d.stock_qty, d.image_url
FROM (VALUES
    ('Milk',       'दूध',       'Tonned milk, dairy farm fresh', 'टोंड दूध',         30.00, 'L',    400, 'https://images.unsplash.com/photo-1550583724-b2692b85b150?w=400'),
    ('Curd',       'दही',       'Fresh set curd',                'ताजा जमा दही',      45.00, 'kg',   150, 'https://images.unsplash.com/photo-1553545985-1e0d8781d2db?w=400'),
    ('Paneer',     'पनीर',      'Malai paneer',                  'मलाई पनीर',        160.00, 'kg',   90,  'https://images.unsplash.com/photo-1631452180519-c014fe946bc7?w=400'),
    ('Eggs',       'अंडे',      'Farm eggs',                     'फार्म अंडे',         7.50, 'piece',500, 'https://images.unsplash.com/photo-1506976785307-8732e854ad03?w=400')
) AS d(name_en, name_hi, description_en, description_hi, price, unit, stock_qty, image_url)
WHERE NOT EXISTS (SELECT 1 FROM products p WHERE p.name_en = d.name_en AND p.category_id = (SELECT id FROM categories WHERE name_en = 'Dairy & Eggs'));

-- Bakery & Snacks
INSERT INTO products (category_id, name_en, name_hi, description_en, description_hi, price, unit, stock_qty, image_url)
SELECT (SELECT id FROM categories WHERE name_en = 'Bakery & Snacks'),
       d.name_en, d.name_hi, d.description_en, d.description_hi, d.price, d.unit, d.stock_qty, d.image_url
FROM (VALUES
    ('White Bread',  'सफेद ब्रेड', 'Soft sliced bread', 'नरम स्लाइस ब्रेड',  35.00, 'piece', 120, 'https://images.unsplash.com/photo-1509440159596-0249088772ff?w=400'),
    ('Biscuits',     'बिस्कुट',    'Cream biscuits',    'क्रीम बिस्कुट',       20.00, 'pack',  200, 'https://images.unsplash.com/photo-1558961363-fa8fdf82db35?w=400'),
    ('Namkeen',      'नमकीन',      'Crispy namkeen mix', 'कुरकुरी नमकीन',      70.00, 'pack',  130, 'https://images.unsplash.com/photo-1626082927389-6cd097ce82a6?w=400')
) AS d(name_en, name_hi, description_en, description_hi, price, unit, stock_qty, image_url)
WHERE NOT EXISTS (SELECT 1 FROM products p WHERE p.name_en = d.name_en AND p.category_id = (SELECT id FROM categories WHERE name_en = 'Bakery & Snacks'));

-- Beverages
INSERT INTO products (category_id, name_en, name_hi, description_en, description_hi, price, unit, stock_qty, image_url)
SELECT (SELECT id FROM categories WHERE name_en = 'Beverages'),
       d.name_en, d.name_hi, d.description_en, d.description_hi, d.price, d.unit, d.stock_qty, d.image_url
FROM (VALUES
    ('Tea Leaves',   'चाय पत्ती', 'Premium tea leaves',  'बेहतरीन चाय पत्ती',       95.00, 'pack', 120, 'https://images.unsplash.com/photo-1564890369478-c89ca6d9cde9?w=400'),
    ('Coffee',       'कॉफी',     'Instant coffee',       'इंस्टेंट कॉफी',          180.00, 'pack', 90,  'https://images.unsplash.com/photo-1514432324607-a09d9b4aefdd?w=400'),
    ('Mineral Water','मिनरल वाटर','Packaged drinking water','पैक्ड पेयजल',        20.00, 'L',    300, 'https://images.unsplash.com/photo-1548839140-29a749e1cf4d?w=400')
) AS d(name_en, name_hi, description_en, description_hi, price, unit, stock_qty, image_url)
WHERE NOT EXISTS (SELECT 1 FROM products p WHERE p.name_en = d.name_en AND p.category_id = (SELECT id FROM categories WHERE name_en = 'Beverages'));

-- Personal Care
INSERT INTO products (category_id, name_en, name_hi, description_en, description_hi, price, unit, stock_qty, image_url)
SELECT (SELECT id FROM categories WHERE name_en = 'Personal Care'),
       d.name_en, d.name_hi, d.description_en, d.description_hi, d.price, d.unit, d.stock_qty, d.image_url
FROM (VALUES
    ('Bathing Soap',  'नहाने का साबुन', 'Gentle bathing soap',      'कोमल नहाने का साबुन',    40.00, 'piece', 200, 'https://images.unsplash.com/photo-1607006317206-bdb6d9f7e6a4?w=400'),
    ('Shampoo',       'शैम्पू',        'Hair care shampoo',         'हेयर केयर शैम्पू',       99.00, 'piece', 150, 'https://images.unsplash.com/photo-1556228720-195a672e8a03?w=400'),
    ('Toothpaste',    'टूथपेस्ट',       'Fresh mint toothpaste',     'फ्रेश मिंट टूथपेस्ट',   55.00, 'piece', 180, 'https://images.unsplash.com/photo-1583947581924-860bda6a26df?w=400')
) AS d(name_en, name_hi, description_en, description_hi, price, unit, stock_qty, image_url)
WHERE NOT EXISTS (SELECT 1 FROM products p WHERE p.name_en = d.name_en AND p.category_id = (SELECT id FROM categories WHERE name_en = 'Personal Care'));

-- Household Essentials
INSERT INTO products (category_id, name_en, name_hi, description_en, description_hi, price, unit, stock_qty, image_url)
SELECT (SELECT id FROM categories WHERE name_en = 'Household Essentials'),
       d.name_en, d.name_hi, d.description_en, d.description_hi, d.price, d.unit, d.stock_qty, d.image_url
FROM (VALUES
    ('Detergent',    'डिटर्जेंट',   'Washing detergent powder',   'कपड़े धोने का डिटर्जेंट',  120.00, 'kg',   110, 'https://images.unsplash.com/photo-1584568694244-14fbdf83bd30?w=400'),
    ('Dishwash Bar', 'डिशवाॅश बार', 'Dishwashing soap bar',       'बर्तन धोने का साबुन',      15.00, 'piece', 300, 'https://images.unsplash.com/photo-1585421514738-01798e348b17?w=400'),
    ('Floor Cleaner', 'फ्लोर क्लीनर','Liquid floor cleaner',      'लिक्विड फ्लोर क्लीनर',    90.00, 'L',    130, 'https://images.unsplash.com/photo-1584820927498-cfe5211fd8bf?w=400')
) AS d(name_en, name_hi, description_en, description_hi, price, unit, stock_qty, image_url)
WHERE NOT EXISTS (SELECT 1 FROM products p WHERE p.name_en = d.name_en AND p.category_id = (SELECT id FROM categories WHERE name_en = 'Household Essentials'));

-- Update existing legacy Rice product with an image
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1586201375761-83865001e31c?w=400', description_en = 'Premium rice', description_hi = 'बेहतरीन चावल', unit = 'kg' WHERE name_en = 'Rice';