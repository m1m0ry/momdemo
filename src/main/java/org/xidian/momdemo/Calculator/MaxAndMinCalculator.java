package org.xidian.momdemo.Calculator;

import org.apache.activemq.ActiveMQConnectionFactory;
import javax.jms.*;

public class MaxAndMinCalculator implements MessageListener {

    private static String brokerURL = "tcp://localhost:61616";
    private static ConnectionFactory factory;
    private Connection connection;
    private Session session;
    private MessageProducer producer;
    private Topic topic;

    public static Double max;
    public static Double min;

    public MaxAndMinCalculator(String topicName) throws JMSException {
        super();
        min = null;
        max = null;

        factory = new ActiveMQConnectionFactory(brokerURL);
        connection = factory.createConnection();

        session = connection.createSession(false, Session.AUTO_ACKNOWLEDGE);
        topic = session.createTopic(topicName);
        producer = session.createProducer(topic);

        connection.start();
    }

    public void close() throws JMSException {
        if (connection != null) {
            connection.close();
        }
    }

    protected void finalize() throws Throwable {
        super.finalize();
        close();
    }

    public void onMessage(Message message) {
        try {
            Double number = (Double) ((ObjectMessage) message).getObject();
            if (max == null) {
                max = number;
                min = number;
            }
            if (number > max)
                max = number;
            if (number < min)
                min = number;
        } catch (Exception e) {
            e.printStackTrace();
        }

        try {
            sendMessage();
        } catch (JMSException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        }
    }

    public void sendMessage() throws JMSException {
        Message message = session.createMapMessage();
        message.setDoubleProperty("max", max);
        message.setDoubleProperty("min", min);
        producer.send(message);
        System.out.println("Sent a message!");
    }
}