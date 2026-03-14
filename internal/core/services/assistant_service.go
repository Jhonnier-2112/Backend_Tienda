package services

import (
	"strings"
)

type assistantService struct {
	knowledgeBase map[string]string
}

func NewAssistantService() *assistantService {
	kb := map[string]string{
		"envio":         "Realizamos envíos a todo el país. El tiempo estimado de entrega es de 3 a 5 días hábiles. El costo depende de tu ubicación, pero ofrecemos envío gratuito en compras superiores a $200.000.",
		"entrega":       "Realizamos envíos a todo el país. El tiempo estimado de entrega es de 3 a 5 días hábiles. El costo depende de tu ubicación, pero ofrecemos envío gratuito en compras superiores a $200.000.",
		"devolucion":    "Nuestra política de devoluciones permite cambios hasta 30 días después de la compra, siempre que el producto esté en su estado original y con etiquetas.",
		"cambio":        "Nuestra política de devoluciones permite cambios hasta 30 días después de la compra, siempre que el producto esté en su estado original y con etiquetas.",
		"horario":       "Estamos abiertos de lunes a viernes de 8:00 AM a 6:00 PM, y los sábados de 9:00 AM a 2:00 PM. Nuestra tienda virtual está disponible 24/7.",
		"pago":          "Aceptamos tarjetas de crédito, débito, transferencias bancarias y pagos a través de PSE.",
		"contacto":      "Puedes contactarnos por WhatsApp al +57 300 000 0000 o enviarnos un correo a soporte@virtualstore.com.",
		"quienes somos": "Somos VirtualStore, una tienda dedicada a ofrecer los mejores productos tecnológicos con la más alta calidad y garantía oficial.",
		"productos":     "Contamos con una amplia variedad de productos tecnológicos, desde laptops y smartphones hasta accesorios para el hogar inteligente. ¡Explora nuestro catálogo en la sección de Tienda!",
		"ubicacion":     "Nuestra sede principal está en Bogotá, pero operamos principalmente como tienda virtual para llegar a cada rincón del país.",
		"garantia":      "Todos nuestros productos cuentan con garantía oficial que varía entre 6 y 12 meses dependiendo del fabricante.",
	}

	return &assistantService{
		knowledgeBase: kb,
	}
}

func (s *assistantService) GetResponse(message string) (string, error) {
	message = strings.ToLower(message)

	// Simple keyword matching
	for key, response := range s.knowledgeBase {
		if strings.Contains(message, key) {
			return response, nil
		}
	}

	return "Lo siento, no tengo una respuesta específica para eso en este momento. Mi especialidad son dudas sobre envíos, pagos, devoluciones, horarios e información general de la tienda. ¿Puedo ayudarte con algo de eso?", nil
}
