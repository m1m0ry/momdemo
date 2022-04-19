package org.xidian.momdemo.Calculator;

import org.apache.activemq.ActiveMQConnectionFactory;
import javax.jms.*;

public class MyListener {

    public static void main(String[] args) throws JMSException {
        String brokerURL = "tcp://localhost:61616";
        ConnectionFactory factory = null;
        Connection connection = null;
        Session session = null;
        Topic topic = null;
        MessageConsumer varianceConsumer = null;
        MessageConsumer averageConsumer = null;
        MessageConsumer maxandminConsumer = null;
        VarianceCalculator averageCalculator = null;
        VarianceCalculator varianceCalculator = null;
        MaxAndMinCalculator maxandmininCalculator =null;


        
            factory = new ActiveMQConnectionFactory(brokerURL);
            connection = factory.createConnection();
            session = connection.createSession(false, Session.AUTO_ACKNOWLEDGE);
            topic = session.createTopic("RandomSignal");

            varianceConsumer = session.createConsumer(topic);
            averageConsumer = session.createConsumer(topic);
            maxandminConsumer = session.createConsumer(topic);

            varianceCalculator = new VarianceCalculator("Variance",1000);
            averageCalculator = new VarianceCalculator("Average",1000);
            maxandmininCalculator = new MaxAndMinCalculator("MaxAndMin");

            averageConsumer.setMessageListener(averageCalculator);
            varianceConsumer.setMessageListener(varianceCalculator);
            maxandminConsumer.setMessageListener(maxandmininCalculator);

            connection.start();
    }

}