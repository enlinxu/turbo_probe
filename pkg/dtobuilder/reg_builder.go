package dtobuilder

import (
	"turbo_probe/pkg/proto"
	"math"
	"fmt"
)

func BuildAccountDefEntry(name, displayName, description, verificationRegex string,
	mandatory bool, isSecret bool) *proto.AccountDefEntry {

	fieldType := &proto.CustomAccountDefEntry_PrimitiveValue_{
		PrimitiveValue: proto.CustomAccountDefEntry_STRING,
	}
	entry := &proto.CustomAccountDefEntry{
		Name:              &name,
		DisplayName:       &displayName,
		Description:       &description,
		VerificationRegex: &verificationRegex,
		IsSecret:          &isSecret,
		FieldType:         fieldType,
	}

	customDef := &proto.AccountDefEntry_CustomDefinition{
		CustomDefinition: entry,
	}

	accountDefEntry := &proto.AccountDefEntry{
		Mandatory:  &mandatory,
		Definition: customDef,
	}

	return accountDefEntry
}


type TemplateDTOBuilder struct {
	templateClass              *proto.EntityDTO_EntityType
	templateType               *proto.TemplateDTO_TemplateType
	priority                   *int32
	commoditiesSold            []*proto.TemplateCommodity
	providerCommodityBoughtMap map[*proto.Provider][]*proto.TemplateCommodity
	externalLinks              []*proto.TemplateDTO_ExternalEntityLinkProp

	currentProvider *proto.Provider

	err error
}

// Create a new SupplyChainNode Builder.
// All the new supply chain node are default to use the base template type and priority 0.
func NewTemplateDTOBuilder(entityType proto.EntityDTO_EntityType) *TemplateDTOBuilder {
	templateType := proto.TemplateDTO_BASE
	priority := int32(0)
	return &TemplateDTOBuilder{
		templateClass: &entityType,
		templateType:  &templateType,
		priority:      &priority,
	}
}

// Create a SupplyChainNode
func (scnb *TemplateDTOBuilder) Create() (*proto.TemplateDTO, error) {
	if scnb.err != nil {
		return nil, fmt.Errorf("Cannot create supply chain node because of error: %v", scnb.err)
	}
	return &proto.TemplateDTO{
		TemplateClass:    scnb.templateClass,
		TemplateType:     scnb.templateType,
		TemplatePriority: scnb.priority,
		CommoditySold:    scnb.commoditiesSold,
		CommodityBought:  buildCommodityBought(scnb.providerCommodityBoughtMap),
		ExternalLink:     scnb.externalLinks,
	}, nil
}

// The very basic selling method. If want others, use other names
func (scnb *TemplateDTOBuilder) Sells(templateComm *proto.TemplateCommodity) *TemplateDTOBuilder {
	if scnb.err != nil {
		return scnb
	}

	if scnb.commoditiesSold == nil {
		scnb.commoditiesSold = []*proto.TemplateCommodity{}
	}
	scnb.commoditiesSold = append(scnb.commoditiesSold, templateComm)
	return scnb
}

// set the provider of the SupplyChainNode
func (scnb *TemplateDTOBuilder) Provider(provider proto.EntityDTO_EntityType, pType proto.Provider_ProviderType) *TemplateDTOBuilder {
	if scnb.err != nil {
		return scnb
	}

	if pType == proto.Provider_LAYERED_OVER {
		// TODO, need a separate class to build provider.
		maxCardinality := int32(math.MaxInt32)
		minCardinality := int32(0)
		scnb.currentProvider = &proto.Provider{
			TemplateClass:  &provider,
			ProviderType:   &pType,
			CardinalityMax: &maxCardinality,
			CardinalityMin: &minCardinality,
		}
	} else {
		hostCardinality := int32(1)
		scnb.currentProvider = &proto.Provider{
			TemplateClass:  &provider,
			ProviderType:   &pType,
			CardinalityMax: &hostCardinality,
			CardinalityMin: &hostCardinality,
		}
	}

	return scnb
}

// Add a commodity this node buys from the current provider. The provider must already been specified.
// If there is no provider for this node, does not add the commodity.
func (scnb *TemplateDTOBuilder) Buys(templateComm *proto.TemplateCommodity) *TemplateDTOBuilder {
	if scnb.err != nil {
		return scnb
	}

	if scnb.currentProvider == nil {
		scnb.err = fmt.Errorf("Provider must be set before calling Buys().")
		return scnb
	}

	if scnb.providerCommodityBoughtMap == nil {
		scnb.providerCommodityBoughtMap = make(map[*proto.Provider][]*proto.TemplateCommodity)
	}

	templateCommoditiesSoldByCurrentProvider, exist := scnb.providerCommodityBoughtMap[scnb.currentProvider]
	if !exist {
		templateCommoditiesSoldByCurrentProvider = []*proto.TemplateCommodity{}
	}
	templateCommoditiesSoldByCurrentProvider = append(templateCommoditiesSoldByCurrentProvider, templateComm)
	scnb.providerCommodityBoughtMap[scnb.currentProvider] = templateCommoditiesSoldByCurrentProvider

	return scnb
}

func buildCommodityBought(providerCommodityBoughtMap map[*proto.Provider][]*proto.TemplateCommodity) []*proto.TemplateDTO_CommBoughtProviderProp {
	if len(providerCommodityBoughtMap) == 0 {
		return nil
	}
	commBought := []*proto.TemplateDTO_CommBoughtProviderProp{}
	for provider, templateCommodities := range providerCommodityBoughtMap {
		commBought = append(commBought, &proto.TemplateDTO_CommBoughtProviderProp{
			Key:   provider,
			Value: templateCommodities,
		})
	}
	return commBought
}