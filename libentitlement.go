package libentitlement

import (
	"github.com/docker/libentitlement/entitlement"
	"github.com/docker/libentitlement/context"
	"fmt"
	"reflect"
)

type EntitlementsManager struct {
	context context.Context
	entitlementList []entitlement.Entitlement
}

// NewEntitlementsManager() instantiates an EntitlementsManager object with the given context
// default
func NewEntitlementsManager(ctx context.Context) (*EntitlementsManager) {
	return &EntitlementsManager{context: ctx, entitlementList: make([]entitlement.Entitlement, 0)}
}

func isValidEntitlement(ent entitlement.Entitlement) (bool, error) {
	_, err := ent.Identifier()
	if err != nil {
		return false, err
	}

	_, err = ent.Domain()
	if err != nil {
		return false, err
	}

	_, err = ent.Value()
	if err != nil {
		return false, err
	}

	return true, nil
}

// Add() adds the given entitlements to the current entitlements list and enforce them
func (m *EntitlementsManager) Add(entitlements ...entitlement.Entitlement) error {
	for _, ent := range entitlements {
		if isValid, err := isValidEntitlement(ent); isValid == false {
			return fmt.Errorf("Couldn't add invalid entitlement: %v", err)
		}

		ctx, err := ent.Enforce(m.context)
		if err != nil {
			return err
		}

		m.context = ctx

		m.entitlementList = append(m.entitlementList, ent)
	}

	return nil
}

// HasEntitlement() returns wether the given entitlement is registered in the current entitlements list
func (m *EntitlementsManager) HasEntitlement(ent entitlement.Entitlement) (bool, error) {
	if isValid, err := isValidEntitlement(ent); isValid == false {
		return false, fmt.Errorf("Couldn't validate invalid entitlement: %v", err)
	}

	for _, currEnts := range m.entitlementList {
		// FIXME: we should not check equality on the function callbacks, only the type, domain, id and value
		if reflect.DeepEqual(currEnts, entitlement) {
			return true, nil
		}
	}

	return false, nil
}

// Enforce() applies the constraints on the security context and updates it to be used for the container
func (m *EntitlementsManager) Enforce() error {
	for _, ent := range m.entitlementList {
		if isValid, err := isValidEntitlement(ent); isValid == false {
			return fmt.Errorf("Couldn't enforce invalid entitlement: %v", err)
		}

		// Try to enforce the entitlement on the security context
		newContext, err := ent.Enforce(m.GetContext())
		if err != nil {
			return err
		}

		m.context = newContext
	}

	return nil
}

// GetContext() returns the current state of the security context
func (m *EntitlementsManager) GetContext() (context.Context) {
	return m.context
}