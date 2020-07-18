package cluster

import (
	context "context"

	"github.com/filanov/bm-inventory/models"

	"github.com/filanov/bm-inventory/internal/common"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("insufficient_state", func() {
	var (
		ctx          = context.Background()
		state        API
		db           *gorm.DB
		currentState = clusterStatusInsufficient
		id           strfmt.UUID
		updateReply  *UpdateReply
		updateErr    error
		cluster      common.Cluster
		dbName       = "cluster_insufficient_state"
	)

	BeforeEach(func() {
		db = common.PrepareTestDB(dbName)
		state = &Manager{insufficient: NewInsufficientState(getTestLog(), db)}
		registerManager := NewRegistrar(getTestLog(), db)

		id = strfmt.UUID(uuid.New().String())
		cluster = common.Cluster{Cluster: models.Cluster{
			ID:     &id,
			Status: swag.String(currentState),
		}}

		replyErr := registerManager.RegisterCluster(ctx, &cluster)
		Expect(replyErr).Should(BeNil())
		Expect(swag.StringValue(cluster.Status)).Should(Equal(clusterStatusInsufficient))
		c := geCluster(*cluster.ID, db)
		Expect(swag.StringValue(c.Status)).Should(Equal(clusterStatusInsufficient))
	})

	Context("refresh_state", func() {
		It("not answering requirement to be ready", func() {
			updateReply, updateErr = state.RefreshStatus(ctx, &cluster, db)
			Expect(updateErr).Should(BeNil())
			Expect(updateReply.State).Should(Equal(clusterStatusInsufficient))
			c := geCluster(*cluster.ID, db)
			Expect(swag.StringValue(c.Status)).Should(Equal(clusterStatusInsufficient))
		})

		It("resetting when host in reboot stage", func() {
			addHost(models.HostRoleMaster, models.HostStatusResetting, *cluster.ID, db)
			c := geCluster(*cluster.ID, db)
			Expect(len(c.Hosts)).Should(Equal(1))
			updateHostProgress(c.Hosts[0], models.HostStageRebooting, "rebooting", db)
			updateReply, updateErr = state.RefreshStatus(ctx, &cluster, db)
			Expect(updateErr).Should(BeNil())
			Expect(updateReply.State).Should(Equal(clusterStatusInsufficient))
			c = geCluster(*cluster.ID, db)
			Expect(c.Hosts[0].Progress.CurrentStage).Should(Equal(models.HostStageRebooting))
			Expect(swag.StringValue(c.Hosts[0].Status)).Should(Equal(models.HostStatusResettingPendingUserAction))
		})

		It("answering requirement to be ready", func() {
			addInstallationRequirements(id, db)
			updateReply, updateErr = state.RefreshStatus(ctx, &cluster, db)
			Expect(updateErr).Should(BeNil())
			Expect(updateReply.State).Should(Equal(clusterStatusReady))
			c := geCluster(*cluster.ID, db)
			Expect(swag.StringValue(c.Status)).Should(Equal(clusterStatusReady))

		})
	})

	AfterEach(func() {
		common.DeleteTestDB(db, dbName)
		updateReply = nil
		updateErr = nil
	})
})
