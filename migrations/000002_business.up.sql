CREATE TYPE attachment_type AS ENUM ('image', 'video');

CREATE TABLE business_categories (
  id UUID PRIMARY KEY,
  name VARCHAR(100) UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT 'now()',
  updated_at TIMESTAMP NOT NULL DEFAULT 'now()'
);

CREATE TABLE businesses (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  category_id UUID REFERENCES business_categories(id) ON DELETE CASCADE,
  address TEXT NOT NULL,
  latitude FLOAT,
  longitude FLOAT,
  contact_info JSON,
  hours_of_operation JSON,
  owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT 'now()',
  updated_at TIMESTAMP NOT NULL DEFAULT 'now()'
);

CREATE TABLE business_attachment (
  id UUID PRIMARY KEY,
  business_id UUID NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
  filepath VARCHAR(255) NOT NULL,
  content_type attachment_type NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
