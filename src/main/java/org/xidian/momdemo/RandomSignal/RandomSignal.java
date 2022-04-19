package org.xidian.momdemo.RandomSignal;

import org.apache.activemq.ActiveMQConnectionFactory;
import javax.jms.*;
import java.util.Random;

public class RandomSignal {

    private static String brokerURL = "tcp://localhost:61616";
    private static ConnectionFactory factory;
    private Connection connection;
    private Session session;
    private MessageProducer producer;
    private Topic topic;

    public RandomSignal(String topicName) throws JMSException {

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

    public static void main(String[] args) throws JMSException, InterruptedException {
        RandomSignal generator = new RandomSignal("RandomSignal");
        generator.sendMessage(20,10);
        generator.close();
    }

    public void sendMessage(double a, double b) throws JMSException, InterruptedException {

        Random random = new Random();
        while (true){
            //生成均值为a，方差为b的随机数
            Double sig =  Math.sqrt(b)*random.nextGaussian()+a;
            Message sigMessage = session.createObjectMessage(sig);
            producer.send(sigMessage);
            System.out.println(sig);
            Thread.sleep(100);
        }
    }
}