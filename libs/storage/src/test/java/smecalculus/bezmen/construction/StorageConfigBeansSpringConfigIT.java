package smecalculus.bezmen.construction;

import org.springframework.test.context.TestPropertySource;

@TestPropertySource(properties = {"solution.config.mapping.mode=spring_config"})
class StorageConfigBeansSpringConfigIT extends StorageConfigBeansIT {}
