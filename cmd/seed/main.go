package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"asona/config"
)

const (
	defaultSeed     = 42
	defaultPassword = "password123"
)

func main() {
	log.Println("🌱 Starting Asona database seeder...")

	// Load environment configuration
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "dev"
	}
	config.Init(appEnv)
	cfg := config.GetConfig()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to open database: %v", err)
	}
	defer func() { _ = db.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}
	log.Println("✅ Connected to database successfully")

	// Initialize faker with deterministic seed
	fake := gofakeit.New(defaultSeed)

	// Hash default password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("❌ Failed to hash password: %v", err)
	}

	// Run seeding in a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("❌ Failed to begin transaction: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			log.Fatalf("❌ Seeding panicked: %v", r)
		}
	}()

	// Truncate all tables (reverse dependency order)
	log.Println("🗑️  Truncating all tables...")
	truncateAll(ctx, tx)

	// Seed in FK dependency order
	seedUsers(ctx, tx, fake, string(hashedPassword))
	seedAuthProviders(ctx, tx, fake)
	seedVerificationCodes(ctx, tx, fake)
	seedOrganizations(ctx, tx, fake)
	seedOrganizationMembers(ctx, tx)
	seedWorkplaces(ctx, tx, fake)
	seedWorkplaceMembers(ctx, tx)
	seedWorkplaceInvitations(ctx, tx, fake)
	seedFriendships(ctx, tx)
	seedProjects(ctx, tx, fake)
	seedProjectMembers(ctx, tx)
	seedProjectStatuses(ctx, tx)
	seedTasks(ctx, tx, fake)
	seedTaskAssignees(ctx, tx)
	seedLabels(ctx, tx, fake)
	seedTaskLabels(ctx, tx)
	seedTaskComments(ctx, tx, fake)
	seedTaskLinks(ctx, tx)
	seedTaskWatchers(ctx, tx)
	seedTaskActivityLogs(ctx, tx, fake)
	seedChannels(ctx, tx, fake)
	seedChannelMembers(ctx, tx)
	seedMessages(ctx, tx, fake)
	seedMessageReactions(ctx, tx)

	if err := tx.Commit(); err != nil {
		log.Fatalf("❌ Failed to commit transaction: %v", err)
	}

	log.Println("✅ Seeding completed successfully!")
	log.Printf("📝 Default login credentials: email=user1@example.com, password=%s", defaultPassword)
}

// truncateAll removes all data from tables in reverse FK dependency order.
func truncateAll(ctx context.Context, tx *sql.Tx) {
	tables := []string{
		"message_reactions",
		"message_attachments",
		"messages",
		"channel_members",
		"channels",
		"task_activity_logs",
		"task_watchers",
		"task_links",
		"task_labels",
		"labels",
		"task_comments",
		"task_attachments",
		"task_assignees",
		"tasks",
		"project_statuses",
		"project_members",
		"projects",
		"friendships",
		"workplace_invitations",
		"workplace_members",
		"workplaces",
		"organization_members",
		"organizations",
		"verification_codes",
		"auth_providers",
		"sessions",
		"users",
	}
	query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", strings.Join(tables, ", "))
	mustExec(ctx, tx, query)
	log.Println("   ✅ All tables truncated")
}

// ─── USERS ───────────────────────────────────────────────────────────────────

func seedUsers(ctx context.Context, tx *sql.Tx, fake *gofakeit.Faker, hashedPassword string) {
	log.Println("👤 Seeding users...")
	users := []struct {
		name, username, displayName, email, avatarURL string
		onboardingStatus                              string
		onboardingStep                                int
		completed                                     bool
	}{
		{"Nguyen Van A", "nguyenvana", "Van A", "user1@example.com", "https://i.pravatar.cc/200?u=1", "completed", 3, true},
		{"Tran Thi B", "tranthib", "Thi B", "user2@example.com", "https://i.pravatar.cc/200?u=2", "in_progress", 1, false},
		{"Le Van C", "levanc", "Van C", "user3@example.com", "https://i.pravatar.cc/200?u=3", "pending", 0, false},
		{"Pham Thi D", "phamthid", "Thi D", "user4@example.com", "https://i.pravatar.cc/200?u=4", "in_progress", 2, false},
		{"Hoang Van E", "hoangvane", "Van E", "user5@example.com", "https://i.pravatar.cc/200?u=5", "completed", 3, true},
	}

	for _, u := range users {
		var onboardedAt interface{}
		if u.completed {
			onboardedAt = time.Now().Add(-24 * time.Hour)
		}

		mustExec(ctx, tx, `
			INSERT INTO users (
				name, username, display_name, email, email_verified, password, avatar_url, is_active,
				onboarding_status, onboarding_step, onboarded_at
			)
			VALUES ($1, $2, $3, $4, NOW(), $5, $6, true, $7, $8, $9)`,
			u.name, u.username, u.displayName, u.email, hashedPassword, u.avatarURL, u.onboardingStatus, u.onboardingStep, onboardedAt,
		)
	}
	log.Printf("   ✅ %d users seeded", len(users))
}

// ─── AUTH PROVIDERS ──────────────────────────────────────────────────────────

func seedAuthProviders(ctx context.Context, tx *sql.Tx, fake *gofakeit.Faker) {
	log.Println("🔐 Seeding auth providers...")
	providers := []struct {
		userID            int
		provider          string
		providerAccountID string
	}{
		{1, "google", fake.UUID()},
		{3, "google", fake.UUID()},
	}

	for _, p := range providers {
		mustExec(ctx, tx, `
			INSERT INTO auth_providers (user_id, provider, provider_account_id, access_token, scope)
			VALUES ($1, $2, $3, $4, 'openid email profile')`,
			p.userID, p.provider, p.providerAccountID, fake.UUID(),
		)
	}
	log.Printf("   ✅ %d auth providers seeded", len(providers))
}

// ─── VERIFICATION CODES ─────────────────────────────────────────────────────

func seedVerificationCodes(ctx context.Context, tx *sql.Tx, fake *gofakeit.Faker) {
	log.Println("📧 Seeding verification codes...")
	codes := []struct {
		userID     int
		identifier string
		codeType   string
	}{
		{2, "user2@example.com", "email_verification"},
		{4, "user4@example.com", "email_verification"},
	}

	for _, c := range codes {
		mustExec(ctx, tx, `
			INSERT INTO verification_codes (user_id, identifier, code, type, expires_at)
			VALUES ($1, $2, $3, $4, NOW() + INTERVAL '24 hours')`,
			c.userID, c.identifier, fake.DigitN(6), c.codeType,
		)
	}
	log.Printf("   ✅ %d verification codes seeded", len(codes))
}

// ─── ORGANIZATIONS ──────────────────────────────────────────────────────────

func seedOrganizations(ctx context.Context, tx *sql.Tx, fake *gofakeit.Faker) {
	log.Println("🏢 Seeding organizations...")
	orgs := []struct {
		name, slug, description string
		createdBy               int
	}{
		{"Betasoft Mobile", "betasoft-mobile", fake.Sentence(10), 1},
		{"Startup Vietnam", "startup-vietnam", fake.Sentence(10), 2},
	}

	for _, o := range orgs {
		mustExec(ctx, tx, `
			INSERT INTO organizations (name, slug, logo_url, description, created_by)
			VALUES ($1, $2, $3, $4, $5)`,
			o.name, o.slug, "https://picsum.photos/100/100", o.description, o.createdBy,
		)
	}
	log.Printf("   ✅ %d organizations seeded", len(orgs))
}

func seedOrganizationMembers(ctx context.Context, tx *sql.Tx) {
	log.Println("👥 Seeding organization members...")
	members := []struct {
		orgID  int
		userID int
		role   string
	}{
		{1, 1, "admin"},
		{1, 2, "member"},
		{1, 3, "member"},
		{2, 2, "admin"},
		{2, 4, "sub_admin"},
		{2, 5, "member"},
	}

	for _, m := range members {
		mustExec(ctx, tx, `
			INSERT INTO organization_members (organization_id, user_id, role)
			VALUES ($1, $2, $3)`,
			m.orgID, m.userID, m.role,
		)
	}
	log.Printf("   ✅ %d organization members seeded", len(members))
}

// ─── WORKPLACES ─────────────────────────────────────────────────────────────

func seedWorkplaces(ctx context.Context, tx *sql.Tx, fake *gofakeit.Faker) {
	log.Println("🏠 Seeding workplaces...")
	workplaces := []struct {
		name      string
		size      string
		createdBy int
	}{
		{"Engineering Team", "21-50", 1},
		{"Marketing Team", "6-10", 2},
	}

	for _, w := range workplaces {
		mustExec(ctx, tx, `
			INSERT INTO workplaces (name, icon_url, size, created_by)
			VALUES ($1, $2, $3, $4)`,
			w.name, "https://picsum.photos/64/64", w.size, w.createdBy,
		)
	}
	log.Printf("   ✅ %d workplaces seeded", len(workplaces))
}

func seedWorkplaceMembers(ctx context.Context, tx *sql.Tx) {
	log.Println("👥 Seeding workplace members...")
	members := []struct {
		workplaceID int
		userID      int
		role        string
	}{
		{1, 1, "admin"},
		{1, 2, "member"},
		{1, 3, "member"},
		{1, 4, "member"},
		{2, 2, "admin"},
		{2, 5, "member"},
	}

	for _, m := range members {
		mustExec(ctx, tx, `
			INSERT INTO workplace_members (workplace_id, user_id, role)
			VALUES ($1, $2, $3)`,
			m.workplaceID, m.userID, m.role,
		)
	}
	log.Printf("   ✅ %d workplace members seeded", len(members))
}

// ─── WORKPLACE INVITATIONS ──────────────────────────────────────────────────

func seedWorkplaceInvitations(ctx context.Context, tx *sql.Tx, fake *gofakeit.Faker) {
	log.Println("✉️  Seeding workplace invitations...")
	mustExec(ctx, tx, `
		INSERT INTO workplace_invitations (workplace_id, invite_token, created_by, max_uses, use_count, expires_at)
		VALUES ($1, $2, $3, 10, 2, NOW() + INTERVAL '7 days')`,
		1, fake.UUID(), 1,
	)
	log.Println("   ✅ 1 workplace invitation seeded")
}

// ─── FRIENDSHIPS ────────────────────────────────────────────────────────────

func seedFriendships(ctx context.Context, tx *sql.Tx) {
	log.Println("🤝 Seeding friendships...")
	friendships := []struct {
		requesterID int
		receiverID  int
		status      string
	}{
		{1, 2, "accepted"},
		{1, 3, "accepted"},
		{2, 4, "pending"},
		{3, 5, "blocked"},
	}

	for _, f := range friendships {
		mustExec(ctx, tx, `
			INSERT INTO friendships (requester_id, receiver_id, status)
			VALUES ($1, $2, $3)`,
			f.requesterID, f.receiverID, f.status,
		)
	}
	log.Printf("   ✅ %d friendships seeded", len(friendships))
}

// ─── PROJECTS ───────────────────────────────────────────────────────────────

func seedProjects(ctx context.Context, tx *sql.Tx, fake *gofakeit.Faker) {
	log.Println("📁 Seeding projects...")
	projects := []struct {
		workplaceID int
		name        string
		description string
		color       string
		access      string
		createdBy   int
	}{
		{1, "Asona Platform", fake.Sentence(10), "#3B82F6", "private", 1},
		{1, "Mobile App v2", fake.Sentence(10), "#EF4444", "private", 2},
		{2, "Landing Page", fake.Sentence(10), "#10B981", "public", 2},
	}

	for _, p := range projects {
		mustExec(ctx, tx, `
			INSERT INTO projects (workplace_id, name, description, color, access, created_by)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			p.workplaceID, p.name, p.description, p.color, p.access, p.createdBy,
		)
	}
	log.Printf("   ✅ %d projects seeded", len(projects))
}

func seedProjectMembers(ctx context.Context, tx *sql.Tx) {
	log.Println("👥 Seeding project members...")
	members := []struct {
		projectID int
		userID    int
		role      string
	}{
		{1, 1, "owner"},
		{1, 2, "member"},
		{1, 3, "member"},
		{2, 2, "owner"},
		{2, 4, "member"},
		{3, 2, "owner"},
		{3, 5, "member"},
	}

	for _, m := range members {
		mustExec(ctx, tx, `
			INSERT INTO project_members (project_id, user_id, role)
			VALUES ($1, $2, $3)`,
			m.projectID, m.userID, m.role,
		)
	}
	log.Printf("   ✅ %d project members seeded", len(members))
}

func seedProjectStatuses(ctx context.Context, tx *sql.Tx) {
	log.Println("📊 Seeding project statuses...")
	statuses := []struct {
		projectID int
		name      string
		color     string
		position  int
		isDefault bool
	}{
		// Project 1: Asona Platform
		{1, "To Do", "#94A3B8", 0, true},
		{1, "In Progress", "#3B82F6", 1, false},
		{1, "In Review", "#F59E0B", 2, false},
		{1, "Done", "#10B981", 3, false},
		// Project 2: Mobile App v2
		{2, "Backlog", "#94A3B8", 0, true},
		{2, "In Progress", "#3B82F6", 1, false},
		{2, "QA Testing", "#8B5CF6", 2, false},
		{2, "Done", "#10B981", 3, false},
		// Project 3: Landing Page
		{3, "To Do", "#94A3B8", 0, true},
		{3, "In Progress", "#3B82F6", 1, false},
		{3, "Done", "#10B981", 2, false},
	}

	for _, s := range statuses {
		mustExec(ctx, tx, `
			INSERT INTO project_statuses (project_id, name, color, position, is_default)
			VALUES ($1, $2, $3, $4, $5)`,
			s.projectID, s.name, s.color, s.position, s.isDefault,
		)
	}
	log.Printf("   ✅ %d project statuses seeded", len(statuses))
}

// ─── TASKS ──────────────────────────────────────────────────────────────────

func seedTasks(ctx context.Context, tx *sql.Tx, fake *gofakeit.Faker) {
	log.Println("📋 Seeding tasks...")
	tasks := []struct {
		projectID  int
		statusID   int
		title      string
		priority   string
		position   float64
		createdBy  int
		reporterID int
	}{
		// Project 1 tasks (status IDs: 1=To Do, 2=In Progress, 3=In Review, 4=Done)
		{1, 1, "Setup CI/CD pipeline", "high", 1000, 1, 1},
		{1, 2, "Implement user authentication", "highest", 2000, 1, 1},
		{1, 2, "Design database schema", "high", 3000, 2, 1},
		{1, 3, "Write API documentation", "medium", 4000, 3, 2},
		{1, 4, "Project kickoff meeting", "low", 5000, 1, 1},
		// Project 2 tasks (status IDs: 5=Backlog, 6=In Progress, 7=QA, 8=Done)
		{2, 5, "Design splash screen", "medium", 1000, 2, 2},
		{2, 6, "Implement push notifications", "high", 2000, 4, 2},
		{2, 8, "Setup React Native project", "highest", 3000, 2, 2},
		// Project 3 tasks (status IDs: 9=To Do, 10=In Progress, 11=Done)
		{3, 9, "Create hero section", "high", 1000, 2, 2},
		{3, 10, "Add testimonials section", "medium", 2000, 5, 2},
	}

	for _, t := range tasks {
		dueDate := fake.DateRange(time.Now(), time.Now().AddDate(0, 1, 0))
		mustExec(ctx, tx, `
			INSERT INTO tasks (project_id, status_id, title, description, priority, position, due_date, created_by, reporter_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			t.projectID, t.statusID, t.title, fake.Sentence(15), t.priority, t.position, dueDate, t.createdBy, t.reporterID,
		)
	}
	log.Printf("   ✅ %d tasks seeded", len(tasks))
}

func seedTaskAssignees(ctx context.Context, tx *sql.Tx) {
	log.Println("👤 Seeding task assignees...")
	assignees := []struct {
		taskID     int
		userID     int
		assignedBy int
	}{
		{1, 3, 1},
		{2, 1, 1},
		{2, 2, 1},
		{3, 2, 1},
		{4, 3, 2},
		{6, 4, 2},
		{7, 4, 2},
		{9, 5, 2},
		{10, 5, 2},
	}

	for _, a := range assignees {
		mustExec(ctx, tx, `
			INSERT INTO task_assignees (task_id, user_id, assigned_by)
			VALUES ($1, $2, $3)`,
			a.taskID, a.userID, a.assignedBy,
		)
	}
	log.Printf("   ✅ %d task assignees seeded", len(assignees))
}

// ─── LABELS ─────────────────────────────────────────────────────────────────

func seedLabels(ctx context.Context, tx *sql.Tx, fake *gofakeit.Faker) {
	_ = fake // suppress unused warning
	log.Println("🏷️  Seeding labels...")
	labels := []struct {
		projectID int
		name      string
		color     string
	}{
		{1, "bug", "#EF4444"},
		{1, "feature", "#3B82F6"},
		{1, "documentation", "#8B5CF6"},
		{1, "enhancement", "#10B981"},
		{2, "bug", "#EF4444"},
		{2, "ui/ux", "#F59E0B"},
		{3, "design", "#EC4899"},
	}

	for _, l := range labels {
		mustExec(ctx, tx, `
			INSERT INTO labels (project_id, name, color)
			VALUES ($1, $2, $3)`,
			l.projectID, l.name, l.color,
		)
	}
	log.Printf("   ✅ %d labels seeded", len(labels))
}

func seedTaskLabels(ctx context.Context, tx *sql.Tx) {
	log.Println("🔖 Seeding task labels...")
	taskLabels := []struct {
		taskID  int
		labelID int
	}{
		{1, 2}, // Setup CI/CD -> feature
		{2, 2}, // Implement auth -> feature
		{3, 3}, // Design DB schema -> documentation
		{4, 3}, // Write API docs -> documentation
		{6, 6}, // Design splash -> ui/ux
		{9, 7}, // Create hero -> design
	}

	for _, tl := range taskLabels {
		mustExec(ctx, tx, `
			INSERT INTO task_labels (task_id, label_id)
			VALUES ($1, $2)`,
			tl.taskID, tl.labelID,
		)
	}
	log.Printf("   ✅ %d task labels seeded", len(taskLabels))
}

// ─── TASK COMMENTS ──────────────────────────────────────────────────────────

func seedTaskComments(ctx context.Context, tx *sql.Tx, fake *gofakeit.Faker) {
	log.Println("💬 Seeding task comments...")
	comments := []struct {
		taskID   int
		authorID int
	}{
		{1, 1},
		{1, 3},
		{2, 2},
		{2, 1},
		{4, 3},
	}

	for _, c := range comments {
		mustExec(ctx, tx, `
			INSERT INTO task_comments (task_id, author_id, content)
			VALUES ($1, $2, $3)`,
			c.taskID, c.authorID, fake.Sentence(12),
		)
	}
	log.Printf("   ✅ %d task comments seeded", len(comments))
}

// ─── TASK LINKS ─────────────────────────────────────────────────────────────

func seedTaskLinks(ctx context.Context, tx *sql.Tx) {
	log.Println("🔗 Seeding task links...")
	links := []struct {
		sourceID  int
		targetID  int
		linkType  string
		createdBy int
	}{
		{2, 3, "blocks", 1},     // Auth blocks DB schema
		{4, 2, "relates_to", 2}, // API docs relates to Auth
	}

	for _, l := range links {
		mustExec(ctx, tx, `
			INSERT INTO task_links (source_id, target_id, link_type, created_by)
			VALUES ($1, $2, $3, $4)`,
			l.sourceID, l.targetID, l.linkType, l.createdBy,
		)
	}
	log.Printf("   ✅ %d task links seeded", len(links))
}

// ─── TASK WATCHERS ──────────────────────────────────────────────────────────

func seedTaskWatchers(ctx context.Context, tx *sql.Tx) {
	log.Println("👁️  Seeding task watchers...")
	watchers := []struct {
		taskID int
		userID int
	}{
		{1, 1},
		{1, 2},
		{2, 1},
		{2, 3},
	}

	for _, w := range watchers {
		mustExec(ctx, tx, `
			INSERT INTO task_watchers (task_id, user_id)
			VALUES ($1, $2)`,
			w.taskID, w.userID,
		)
	}
	log.Printf("   ✅ %d task watchers seeded", len(watchers))
}

// ─── TASK ACTIVITY LOGS ─────────────────────────────────────────────────────

func seedTaskActivityLogs(ctx context.Context, tx *sql.Tx, fake *gofakeit.Faker) {
	_ = fake // suppress unused warning
	log.Println("📜 Seeding task activity logs...")
	logs := []struct {
		taskID  int
		actorID int
		action  string
	}{
		{1, 1, "created"},
		{2, 1, "created"},
		{2, 1, "status_changed"},
		{5, 1, "completed"},
	}

	for _, l := range logs {
		mustExec(ctx, tx, `
			INSERT INTO task_activity_logs (task_id, actor_id, action, old_value, new_value)
			VALUES ($1, $2, $3, $4, $5)`,
			l.taskID, l.actorID, l.action, nil, nil,
		)
	}
	log.Printf("   ✅ %d task activity logs seeded", len(logs))
}

// ─── CHANNELS ───────────────────────────────────────────────────────────────

func seedChannels(ctx context.Context, tx *sql.Tx, fake *gofakeit.Faker) {
	_ = fake // suppress unused warning
	log.Println("📺 Seeding channels...")
	channels := []struct {
		workplaceID *int
		projectID   *int
		name        string
		channelType string
		createdBy   int
	}{
		{intPtr(1), nil, "general", "global", 1},
		{intPtr(1), intPtr(1), "asona-dev", "project", 1},
		{nil, nil, "", "dm", 2},
	}

	for _, ch := range channels {
		var nameVal interface{} = ch.name
		if ch.name == "" {
			nameVal = nil
		}
		mustExec(ctx, tx, `
			INSERT INTO channels (workplace_id, project_id, name, type, created_by)
			VALUES ($1, $2, $3, $4, $5)`,
			ch.workplaceID, ch.projectID, nameVal, ch.channelType, ch.createdBy,
		)
	}
	log.Printf("   ✅ %d channels seeded", len(channels))
}

func seedChannelMembers(ctx context.Context, tx *sql.Tx) {
	log.Println("👥 Seeding channel members...")
	members := []struct {
		channelID int
		userID    int
	}{
		{1, 1}, {1, 2}, {1, 3}, {1, 4},
		{2, 1}, {2, 2}, {2, 3},
		{3, 2}, {3, 3},
	}

	for _, m := range members {
		mustExec(ctx, tx, `
			INSERT INTO channel_members (channel_id, user_id)
			VALUES ($1, $2)`,
			m.channelID, m.userID,
		)
	}
	log.Printf("   ✅ %d channel members seeded", len(members))
}

// ─── MESSAGES ───────────────────────────────────────────────────────────────

func seedMessages(ctx context.Context, tx *sql.Tx, fake *gofakeit.Faker) {
	log.Println("✉️  Seeding messages...")
	messages := []struct {
		channelID int
		senderID  int
		parentID  *int
	}{
		{1, 1, nil},
		{1, 2, nil},
		{1, 3, nil},
		{1, 1, intPtr(1)}, // Reply to message 1
		{1, 2, intPtr(1)}, // Reply to message 1
		{2, 1, nil},
		{2, 2, nil},
		{2, 3, nil},
		{3, 2, nil},
		{3, 3, nil},
	}

	for _, m := range messages {
		mustExec(ctx, tx, `
			INSERT INTO messages (channel_id, sender_id, parent_id, content)
			VALUES ($1, $2, $3, $4)`,
			m.channelID, m.senderID, m.parentID, fake.Sentence(10),
		)
	}
	log.Printf("   ✅ %d messages seeded", len(messages))
}

func seedMessageReactions(ctx context.Context, tx *sql.Tx) {
	log.Println("😊 Seeding message reactions...")
	reactions := []struct {
		messageID int
		userID    int
		emoji     string
	}{
		{1, 2, "👍"},
		{1, 3, "❤️"},
		{2, 1, "🎉"},
		{6, 2, "👀"},
		{9, 3, "😂"},
	}

	for _, r := range reactions {
		mustExec(ctx, tx, `
			INSERT INTO message_reactions (message_id, user_id, emoji)
			VALUES ($1, $2, $3)`,
			r.messageID, r.userID, r.emoji,
		)
	}
	log.Printf("   ✅ %d message reactions seeded", len(reactions))
}

// ─── HELPERS ────────────────────────────────────────────────────────────────

// mustExec executes a SQL statement and panics on error.
func mustExec(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) {
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		panic(fmt.Sprintf("seed exec failed: %v\nQuery: %s\nArgs: %v", err, query, args))
	}
}

// intPtr returns a pointer to an int value.
func intPtr(v int) *int {
	return &v
}
