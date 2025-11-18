# Supabase Authentication Setup Guide

This guide explains how to set up Supabase authentication with your Go API.

## Prerequisites

1. A Supabase project (create one at [supabase.com](https://supabase.com))
2. Go 1.19+ installed
3. PostgreSQL database (can be Supabase's hosted database)

## Step 1: Configure Environment Variables

Create a `.env` file in your project root with the following variables:

```bash
# Supabase Configuration
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key-here

# Database Configuration
DATABASE_URL=postgresql://username:password@localhost:5432/database_name
DIRECT_URL=postgresql://username:password@localhost:5432/database_name
```

### Getting Supabase Credentials

1. Go to your Supabase project dashboard
2. Navigate to Settings > API
3. Copy the Project URL and anon/public key
4. Update your `.env` file with these values

## Step 2: Update Database Schema

Update your Prisma schema (`prisma/schema.prisma`) to include fields that link to Supabase users if needed. For example:

```prisma
model User {
  id            String   @id @default(uuid())
  supabaseUserId String  @unique @map("supabase_user_id")
  // ... other fields
}
```

## Step 3: Regenerate Prisma Client

After updating the schema, regenerate the Prisma client:

```bash
./regenerate_prisma.sh
```

Or manually:

```bash
go run github.com/steebchen/prisma-client-go generate
```

## Step 4: Run Database Migrations

Apply the schema changes to your database:

```bash
# If using Prisma migrations
npx prisma migrate dev

# Or if using direct SQL
npx prisma db push
```

## How It Works

1. **Authentication**: The `AuthMiddleware` verifies Supabase JWT tokens
2. **User Context**: The authenticated user's Supabase ID is stored in the request context
3. **Data Isolation**: Each authenticated user can only access their own data

## Security Features

- JWT tokens are verified using Supabase's public keys
- User IDs are extracted from verified tokens
- No manual user ID manipulation possible

## Troubleshooting

### Common Issues

1. **JWT verification failed**: Check your Supabase URL and anon key
2. **Database connection errors**: Verify your DATABASE_URL and DIRECT_URL

### Debug Mode

Enable debug logging by setting the gin mode to "debug" in your config:

```yaml
server:
  ginmode: "debug"
```

## Support

If you encounter issues:

1. Check the logs for detailed error messages
2. Verify your environment variables
3. Ensure your Supabase project is properly configured
4. Check that the Prisma client has been regenerated
