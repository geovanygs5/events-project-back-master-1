package main

import (
	"crypto/tls"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"net/smtp"
	"time"
)

// Definimos la estructura de datos de nuestros items
type Item struct {
	SK      string `json:"sk"`
	EventID string `json:"eventid"`
	Email   string `json:"email"`
	Date    string `json:"date"`
}

type Eventitem struct {
	SK   string `json:"sk"`
	Name string `json:"name"`
	Date string `json:"date"`
	Hour string `json:"hour"`
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest() {
	// Creamos una sesión con AWS para acceder a DynamoDB
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	svc := dynamodb.New(sess)

	// Escanear la tabla para obtener todos los items
	/*
		params := &dynamodb.ScanInput{
			TableName: aws.String("eventstable"),
		}

		result, _ := svc.Scan(params)
	*/

	input := &dynamodb.QueryInput{
		TableName: aws.String("eventstable"),
		IndexName: aws.String("sk-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"sk": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String("METADATA#USER"),
					},
				},
			},
		},
	}

	result, err := svc.Query(input)
	//fmt.Printf("%v", result)

	if err != nil {
		fmt.Printf("%v", err)
	}

	// Desenmarcar los resultados en nuestra estructura de datos
	items := []Item{}
	dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	colombiaLocation := time.FixedZone("COT", -5*60*60)
	loc, _ := time.LoadLocation("America/Bogota")

	layout := "2006-01-02 15:04:05"
	currentTime := time.Now().In(colombiaLocation)
	fmt.Printf("La fecha actual es: %v", currentTime)

	// Recorrer cada item y verificar si cumple los criterios
	for _, item := range items {
		if item.SK == "METADATA#USER" {
			eventItem := findEventItem(item.EventID)
			datetime, err := time.ParseInLocation(layout, eventItem.Date+" "+eventItem.Hour+":00", loc)
			if err != nil {
				fmt.Printf("EL error es: %v", err)
				return
			}
			fmt.Printf("%v", datetime)

			duration := currentTime.Sub(datetime)

			fmt.Printf("FALTAN %v HORAS", duration.Hours()*-1)
			if duration.Hours()*-1 <= 24 {
				sendEmail(item.Email)
			}
		}
	}

	fmt.Println("Todo esta ok")
}

func findEventItem(eventID string) Eventitem {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	svc := dynamodb.New(sess)
	// Realizamos una consulta de un solo elemento en DynamoDB
	params := &dynamodb.GetItemInput{
		TableName: aws.String("eventstable"),
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(eventID),
			},
			"sk": {
				S: aws.String("METADATA#EVENT"),
			},
		},
	}

	result, err := svc.GetItem(params)
	if err != nil {
		fmt.Printf("El error de servicio es %v:", err)
		panic("errorunmarshal")
	}

	// Desenmarcamos los resultados en nuestra estructura de datos
	var eventItem Eventitem
	err1 := dynamodbattribute.UnmarshalMap(result.Item, &eventItem)
	if err1 != nil {
		fmt.Printf("El error1 es %v:", err1)
		panic("errorunmarshal")
	}

	//fmt.Printf("Evento = %s", eventItem.Hour)

	return eventItem
}

func sendEmail(email string) {
	/*
		//first try
		fmt.Println("Entro a la funcion senEmail")
		// Nos autenticamos en el servidor SMTP
		auth := smtp.PlainAuth("", "no-responder@events.instanceshape.com", "Prx3PA18,Agm", "mail.events.instanceshape.com")

		fmt.Printf("LA AUTORIZACION ES: %v", auth)

		// Componemos y enviamos el correo electrónico
		err := smtp.SendMail("mail.events.instanceshape.com:465", auth, "no-responder@events.instanceshape.com",
			[]string{email}, []byte("Subject: Recordatorio de evento\n\nTienes un evento mañana, no lo olvides!"))

		// Imprimimos un registro
		if err != nil {
			fmt.Printf("Unable to send email: %s", err)
			return
		}

		fmt.Printf("Email sent to: %s", email)
	*/
	fmt.Println("Entro a la funcion senEmail")
	// Nos autenticamos en el servidor SMTP
	auth := smtp.PlainAuth("", "no-responder@events.instanceshape.com", "Prx3PA18,Agm", "mail.events.instanceshape.com")

	fmt.Printf("LA AUTORIZACION ES: %v", auth)

	// Configuramos TLS
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "mail.events.instanceshape.com",
	}

	// Conectamos al servidor SMTP
	conn, err := tls.Dial("tcp", "mail.events.instanceshape.com:465", tlsconfig)
	if err != nil {
		fmt.Printf("Unable to connect to SMTP server: %s", err)
		return
	}

	// Creamos un cliente SMTP
	client, err := smtp.NewClient(conn, "mail.events.instanceshape.com")
	if err != nil {
		fmt.Printf("Unable to create SMTP client: %s", err)
		return
	}

	// Autenticamos
	if err := client.Auth(auth); err != nil {
		fmt.Printf("Unable to authenticate: %s", err)
		return
	}

	// Componemos y enviamos el correo electrónico
	if err := client.Mail("no-responder@events.instanceshape.com"); err != nil {
		fmt.Printf("Unable to set sender: %s", err)
		return
	}
	if err := client.Rcpt(email); err != nil {
		fmt.Printf("Unable to set recipient: %s", err)
		return
	}
	w, err := client.Data()
	if err != nil {
		fmt.Printf("Unable to create data writer: %s", err)
		return
	}
	_, err = w.Write([]byte("From: no-responder@events.instanceshape.com\r\nSubject: Recordatorio de evento\r\n\r\nTienes un evento mañana, no lo olvides!"))

	//_, err = w.Write([]byte("Subject: Recordatorio de evento\r\n\r\nTienes un evento mañana, no lo olvides!"))
	if err != nil {
		fmt.Printf("Unable to write data: %s", err)
		return
	}
	if err := w.Close(); err != nil {
		fmt.Printf("Unable to close data writer: %s", err)
		return
	}

	// Cerramos la conexión
	if err := client.Quit(); err != nil {
		fmt.Printf("Unable to quit SMTP client: %s", err)
		return
	}

	fmt.Printf("Email sent to: %s", email)
}

//cron(0/5 * * * ? *)
