package smecalculus.bezmen.messaging;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.when;

import java.util.UUID;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.test.context.ContextConfiguration;
import org.springframework.test.context.junit.jupiter.SpringExtension;
import smecalculus.bezmen.construction.SepulkaClientBeans;
import smecalculus.bezmen.core.MessageDm.RegistrationRequest;
import smecalculus.bezmen.core.MessageDmEg;
import smecalculus.bezmen.core.SepulkaService;

@ExtendWith(SpringExtension.class)
@ContextConfiguration(classes = SepulkaClientBeans.class)
abstract class SepulkaClientIT {

    @Autowired
    protected SepulkaClient externalClient;

    @Autowired
    private SepulkaService serviceMock;

    @Test
    void shouldRegisterSepulka() {
        // given
        var externalId = UUID.randomUUID().toString();
        // and
        var request = MessageEmEg.registrationRequest(externalId);
        // and
        when(serviceMock.register(any(RegistrationRequest.class)))
                .thenReturn(MessageDmEg.registrationResponse(externalId).build());
        // and
        var expectedResponse = MessageEmEg.registrationResponse(externalId);
        // when
        var actualResponse = externalClient.register(request);
        // then
        assertThat(actualResponse)
                .usingRecursiveComparison()
                .ignoringExpectedNullFields()
                .isEqualTo(expectedResponse);
    }
}
