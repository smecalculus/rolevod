package smecalculus.rolevod.registration;

import static java.time.Duration.ofSeconds;
import static org.assertj.core.api.Assertions.assertThat;
import static org.awaitility.Awaitility.await;

import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.condition.EnabledIfSystemProperty;
import org.junit.jupiter.api.extension.ExtendWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.test.context.ContextConfiguration;
import org.springframework.test.context.junit.jupiter.SpringExtension;
import smecalculus.rolevod.StandBeans;
import smecalculus.rolevod.messaging.RolevodClient;
import smecalculus.rolevod.messaging.SepulkaMessageEmEg;

@ExtendWith(SpringExtension.class)
@ContextConfiguration(classes = StandBeans.class)
public class SepulkaTest {

    @Autowired
    private RolevodClient rolevodClient;

    @BeforeAll
    void beforeAll() {
        await("isReady").atMost(ofSeconds(5)).until(rolevodClient::isReady);
    }

    @Test
    @Tag("smoke")
    void shouldRegisterSepulka() {
        // given
        var request = SepulkaMessageEmEg.registrationRequest();
        // and
        var expectedResponse = SepulkaMessageEmEg.registrationResponse(request.getExternalId());
        // when
        var actualResponse = rolevodClient.register(request);
        // then
        assertThat(actualResponse).usingRecursiveComparison().isEqualTo(expectedResponse);
    }

    @Test
    @EnabledIfSystemProperty(named = "storage.protocol.mode", matches = "postgres")
    void postgresSpecificTest() {
        // empty
    }

    @Test
    @EnabledIfSystemProperty(named = "storage.protocol.mode", matches = "sqlite")
    void sqliteSpecificTest() {
        // empty
    }
}
