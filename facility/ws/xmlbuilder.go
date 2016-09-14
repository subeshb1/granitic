package ws

import (
	"errors"
	"github.com/graniticio/granitic/config"
	"github.com/graniticio/granitic/instance"
	"github.com/graniticio/granitic/ioc"
	"github.com/graniticio/granitic/logging"
	"github.com/graniticio/granitic/ws"
	"github.com/graniticio/granitic/ws/xml"
)

const (
	xmlResponseWriterName = instance.FrameworkPrefix + "XmlResponseWriter"
	xmlUnmarshallerName   = instance.FrameworkPrefix + "XmlUnmarshaller"
	templateMode          = "TEMPLATE"
	marshalMode           = "MARSHAL"
)

type XMLWsFacilityBuilder struct {
}

func (fb *XMLWsFacilityBuilder) BuildAndRegister(lm *logging.ComponentLoggerManager, ca *config.ConfigAccessor, cc *ioc.ComponentContainer) error {

	wc := BuildAndRegisterWsCommon(lm, ca, cc)

	um := new(xml.StandardXmlUnmarshaller)
	cc.WrapAndAddProto(xmlUnmarshallerName, um)

	mode, _ := ca.StringVal("XmlWs.ResponseMode")

	var rw ws.WsResponseWriter

	switch mode {
	case templateMode:
		rw = fb.createTemplateComponents(ca, cc, wc)
	case marshalMode:
		rw = fb.createMarshalComponents(ca, cc, wc)
	default:
		return errors.New("XmlWs.ResponseMode must be set to either TEMPLATE or MARHSAL")
	}

	BuildRegisterWsDecorator(cc, rw, um, wc, lm)
	OfferAbnormalStatusWriter(rw.(ws.AbnormalStatusWriter), cc, xmlResponseWriterName)

	return nil
}

func (fb *XMLWsFacilityBuilder) createTemplateComponents(ca *config.ConfigAccessor, cc *ioc.ComponentContainer, wc *WsCommon) ws.WsResponseWriter {

	rw := new(xml.TemplatedXMLResponseWriter)
	ca.Populate("XmlWs.ResponseWriter", rw)
	cc.WrapAndAddProto(xmlResponseWriterName, rw)

	rw.FrameworkErrors = wc.FrameworkErrors
	rw.StatusDeterminer = wc.StatusDeterminer

	return rw

}

func (fb *XMLWsFacilityBuilder) createMarshalComponents(ca *config.ConfigAccessor, cc *ioc.ComponentContainer, wc *WsCommon) ws.WsResponseWriter {

	rw := new(ws.MarshallingResponseWriter)
	ca.Populate("XmlWs.ResponseWriter", rw)
	cc.WrapAndAddProto(xmlResponseWriterName, rw)

	rw.StatusDeterminer = wc.StatusDeterminer
	rw.FrameworkErrors = wc.FrameworkErrors

	if !cc.ModifierExists(xmlResponseWriterName, "ErrorFormatter") {
		rw.ErrorFormatter = new(xml.StandardXMLErrorFormatter)
	}

	if !cc.ModifierExists(xmlResponseWriterName, "ResponseWrapper") {
		wrap := new(xml.StandardXMLResponseWrapper)
		rw.ResponseWrapper = wrap
	}

	if !cc.ModifierExists(xmlResponseWriterName, "MarshalingWriter") {

		mw := new(xml.XMLMarshalingWriter)
		ca.Populate("XmlWs.Marshal", mw)
		rw.MarshalingWriter = mw
	}

	return rw

}

func (fb *XMLWsFacilityBuilder) FacilityName() string {
	return "XmlWs"
}

func (fb *XMLWsFacilityBuilder) DependsOnFacilities() []string {
	return []string{}
}