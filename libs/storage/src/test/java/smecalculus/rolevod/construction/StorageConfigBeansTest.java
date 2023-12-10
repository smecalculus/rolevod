package smecalculus.rolevod.construction;

import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;
import static smecalculus.rolevod.configuration.StorageEmEg.storageProps;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import smecalculus.rolevod.configuration.PropsKeeper;
import smecalculus.rolevod.configuration.StorageEm.StorageProps;
import smecalculus.rolevod.configuration.StoragePropsMapper;
import smecalculus.rolevod.configuration.StoragePropsMapperImpl;
import smecalculus.rolevod.validation.EdgeValidator;

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
