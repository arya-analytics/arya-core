package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/route"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	localExec query.Execute
	remote    ServiceRemote
}

func NewService(local query.Execute, remote ServiceRemote) *Service {
	return &Service{remote: remote, localExec: local}
}

func (s *Service) CanHandle(p *query.Pack) bool {
	return catalog().Contains(p.Model().Pointer())
}

func (s *Service) Exec(ctx context.Context, p *query.Pack) error {
	return query.Switch(ctx, p, query.Ops{
		Create:   s.createReplica,
		Retrieve: s.retrieveReplica,
		Delete:   s.deleteReplica,
		Update:   s.updateReplica,
	})
}

// |||| REPLICA ||||

// Abbreviation reminder:
// RR - Range Replica
// CCR - Channel Chunk Replica
const (
	ccrFieldRRID       = "RangeReplicaID"
	ccrFieldRR         = "RangeReplica"
	rrRelNode          = "Node"
	ccrFieldNodeIsHost = "RangeReplica.Node.IsHost"
	ccrRelNode         = "RangeReplica.Node"
	ccrTelemField      = "Telem"
)

// These are the fields we need to make a remote/local decision and send a request.
func nodeFields() []string {
	return []string{"ID", "Address", "IsHost", "RPCPort"}

}

func retrieveRRQuery(m *model.Reflect) query.Query {
	return query.NewRetrieve().
		Model(m.FieldsByName(ccrFieldRR).ToReflect()).
		// These are the fields we need to make a remote/local decision and send a request.
		Relation(rrRelNode, nodeFields()...).
		WherePKs(m.FieldsByName(ccrFieldRRID).ToPKChain())
}

func (s *Service) createReplica(ctx context.Context, p *query.Pack) error {
	// CLARIFICATION: Retrieves information about the rng replicas and nodes model belongs to.
	// It will bind the results to p.Model itself.
	if err := s.localExec(ctx, retrieveRRQuery(p.Model()).Pack()); err != nil {
		return err
	}
	// CLARIFICATION: Now that we have the RangeReplicas.Node.IsHost field populated, we can switch on it.
	return replicaNodeIsHostSwitch(
		p.Model(),
		func(m *model.Reflect) error { return s.localExec(ctx, query.NewCreate().Model(m).Pack()) },
		func(m *model.Reflect) error { return s.remote.Create(ctx, remCreateOpts(m)) },
	)
}

// retrieveRequiredFields returns the minimum set of fields we need to complete a channel chunk replica retrieve
// request. We need this info to resolve the node that the replica belongs to.
func retrieveRequiredFields() []string {
	return []string{"ID", "ChannelChunkID", "RangeReplicaID"}
}

func augmentedRetrieveQuery(p *query.Pack) *query.Pack {
	fldsOpt, _ := query.RetrieveFieldsOpt(p)
	query.NewFieldsOpt(p, fldsOpt.AllExcept(ccrTelemField).Append(retrieveRequiredFields()...)...)
	query.NewRelationOpt(p, "RangeReplica", "ID")
	query.NewRelationOpt(p, "RangeReplica.Node", "ID", "Address", "IsHost", "RPCPort")
	return p
}

func (s *Service) retrieveReplica(ctx context.Context, p *query.Pack) error {
	fldsOpt, fldsOptOk := query.RetrieveFieldsOpt(p)
	// CLARIFICATION: If we don't need to retrieve any telemetry, just
	// run the original query and return the result.
	if fldsOptOk && !fldsOpt.ContainsAny(ccrTelemField) {
		return s.localExec(ctx, p)
	}

	// CLARIFICATION: Retrieves information about the rng replicas and nodes model belongs to.
	// It will bind the results to p.Model itself.
	if err := s.localExec(ctx, augmentedRetrieveQuery(p)); err != nil {
		return err
	}

	// CLARIFICATION: Now that we have the RangeReplicas.Node.IsHost field populated, we can switch on it.
	return replicaNodeIsHostSwitch(
		p.Model(),
		func(m *model.Reflect) error {
			return s.localExec(ctx, query.NewRetrieve().Model(m).WherePKs(m.PKChain()).Pack())
		},
		func(m *model.Reflect) error { return s.remote.Retrieve(ctx, m.Pointer(), remRetrieveOpts(m)) },
	)
}

func preDeleteRetrieveQuery(p *query.Pack) query.Query {
	q := query.NewRetrieve().Model(p.Model().Pointer())
	pkc, pkOk := query.PKOpt(p)
	if pkOk {
		q.WherePKs(pkc.Raw())
	}
	wf, wfOk := query.WhereFieldsOpt(p)
	if wfOk {
		q.WhereFields(wf)
	}
	if !pkOk && !wfOk {
		panic("delete queries require at lease one where")
	}
	q.Relation(ccrFieldRR, "ID").
		Relation(ccrRelNode, nodeFields()...)
	return q
}

func (s *Service) deleteReplica(ctx context.Context, p *query.Pack) error {
	// CLARIFICATION: Retrieves information about the rng replicas and nodes model belongs to.
	// It will bind the results to p .Model itself.
	if err := s.localExec(ctx, preDeleteRetrieveQuery(p).Pack()); err != nil {
		return err
	}
	// CLARIFICATION: Now that we have the RangeReplicas.Node.IsHost field populated, we can switch on it.
	return replicaNodeIsHostSwitch(
		p.Model(),
		func(m *model.Reflect) error {
			return s.localExec(ctx, query.NewDelete().Model(m).WherePKs(m.PKChain()).Pack())
		},
		func(m *model.Reflect) error { return s.remote.Delete(ctx, remDeleteOpts(m)) },
	)
}

func (s *Service) updateReplica(ctx context.Context, p *query.Pack) error {
	if !p.Model().FieldsByName("Telem").AllNonZero() {
		log.
			WithFields(log.Fields{"ID": p.Model().PKChain().Raw()}).
			Warn("can't perform update on replica's telemetry, but was still specified!")
	}
	return s.localExec(ctx, p)
}

// |||| ROUTING ||||

func replicaNodeIsHostSwitch(mRfl *model.Reflect, localF, remoteF func(m *model.Reflect) error) error {
	c := errutil.CatchSimple{}
	route.ModelSwitchBoolean(mRfl,
		ccrFieldNodeIsHost,
		func(_ bool, m *model.Reflect) { c.Exec(func() error { return localF(m) }) },
		func(_ bool, m *model.Reflect) { c.Exec(func() error { return remoteF(m) }) },
	)
	return c.Error()
}

func replicaNodeSwitch(rfl *model.Reflect, action func(node *models.Node, rfl *model.Reflect)) {
	route.ModelSwitchIter(rfl, ccrRelNode, action)
}

/// |||| REMOTE OPTION BUILDING ||||

func remRetrieveOpts(ccr *model.Reflect) (opts []RemoteRetrieveOpts) {
	replicaNodeSwitch(ccr, func(node *models.Node, m *model.Reflect) {
		opts = append(opts, RemoteRetrieveOpts{Node: node, PKC: m.PKChain()})
	})
	return opts

}

func remCreateOpts(ccr *model.Reflect) (opts []RemoteCreateOpts) {
	replicaNodeSwitch(ccr, func(node *models.Node, m *model.Reflect) {
		opts = append(opts, RemoteCreateOpts{Node: node, ChunkReplica: m.Pointer()})
	})
	return opts
}

func remDeleteOpts(ccr *model.Reflect) (opts []RemoteDeleteOpts) {
	replicaNodeSwitch(ccr, func(node *models.Node, m *model.Reflect) {
		opts = append(opts, RemoteDeleteOpts{Node: node, PKC: m.PKChain()})
	})
	return opts
}

// |||| CATALOG ||||

func catalog() model.Catalog {
	return model.Catalog{&models.ChannelChunkReplica{}}
}
