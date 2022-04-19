package org.xidian.momdemo.Calculator;

import org.apache.activemq.ActiveMQConnectionFactory;

import java.util.LinkedList;

import javax.jms.*;

public class AverageCalculator implements MessageListener {

    private static String brokerURL = "tcp://localhost:61616";
    private static ConnectionFactory factory;
    private Connection connection;
    private Session session;
    private MessageProducer producer;
    private Topic topic;

    public static Double mean;
    private static LinkedList<Double> datas ;
    Double sum;
    long nums;
    

    public AverageCalculator(String topicName,long num) throws JMSException {
        super();
        mean = 0.0;
        sum = 0.0;
        nums = num;
        datas=new LinkedList<>();

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
            Double now=(Double) ((ObjectMessage) message).getObject();
            datas.add(now);
            if(datas.size()>nums){
                Double pre=datas.poll();
                //移动均值
                mean=mean+(now-pre)/nums;
            }
            else{
                sum +=now;
                mean=sum/datas.size();
            }
            
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
        Message message = session.createObjectMessage(mean);
		producer.send(message);
		System.out.println("Sent a message!");
    }
}