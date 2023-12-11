package smecalculus.rolevod;

import com.fasterxml.jackson.databind.ObjectMapper;
import java.net.http.HttpClient;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import smecalculus.rolevod.messaging.RolevodClient;
import smecalculus.rolevod.messaging.RolevodClientJavaHttp;

@Configuration(proxyBeanMethods = false)
public class StandBeans {

    @Bean
    RolevodClient rolevodClient() {
        var jsonMapper = new ObjectMapper();
        var client = HttpClient.newHttpClient();
        return new RolevodClientJavaHttp(jsonMapper, client);
    }
}
