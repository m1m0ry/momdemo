package org.xidian.momdemo.Calculator;

import org.apache.activemq.ActiveMQConnectionFactory;

import java.util.LinkedList;

import javax.jms.*;

public class VarianceCalculator implements MessageListener {

    private static String brokerURL = "tcp://localhost:61616";
    private static ConnectionFactory factory;
    private Connection connection;
    private Session session;
    private MessageProducer producer;
    private Topic topic;

    private static LinkedList<Double> datas ;
    public static Double variance;
    private static Double mean;
    private static Double sum;
    private static Double sum2;
    long nums;

    public VarianceCalculator(String topicName,long num) throws JMSException {
        super();
        variance = 0.0;
        sum = 0.0;
        sum2=0.0;
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
                // 移动/滚动 方差
                Double old_mean=mean;
                mean=mean+(now-pre)/nums;
                variance=variance+(mean-old_mean)*((nums-1)*(pre+now)-2*nums*old_mean+2*pre);
            }
            else{
                sum +=now;
                mean=sum/datas.size();
                sum2 += Math.pow(now - mean, 2.0);
                variance = sum2/datas.size();
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
        Message message = session.createObjectMessage(variance);
		producer.send(message);
		System.out.println("Sent a message!");
    }
}