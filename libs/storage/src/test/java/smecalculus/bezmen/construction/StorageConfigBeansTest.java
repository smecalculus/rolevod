package smecalculus.bezmen.construction;

import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;
import static smecalculus.bezmen.configuration.StorageEmEg.storageProps;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import smecalculus.bezmen.configuration.PropsKeeper;
import smecalculus.bezmen.configuration.StorageEm.StorageProps;
import smecalculus.bezmen.configuration.StoragePropsMapper;
import smecalculus.bezmen.configuration.StoragePropsMapperImpl;
import smecalculus.bezmen.validation.EdgeValidator;

class StorageConfigBeansTest {

    private final StorageConfigBeans config = new StorageConfigBeans();
    private final StoragePropsMapper mapper = new StoragePropsMapperImpl();
    private EdgeValidator validatorMock;
    private PropsKeeper keeperMock;

    @BeforeEach
    void setUp() {
        validatorMock = mock(EdgeValidator.class);
        keeperMock = mock(PropsKeeper.class);
    }

    @Test
    void shouldValidateConf() {
        // given
        StorageProps expectedProps = storageProps();
        // and
        when(keeperMock.read("solution.storage", StorageProps.class)).thenReturn(expectedProps);
        // when
        config.storageProps(keeperMock, validatorMock, mapper);
        // then
        verify(validatorMock).validate(expectedProps);
    }
}
