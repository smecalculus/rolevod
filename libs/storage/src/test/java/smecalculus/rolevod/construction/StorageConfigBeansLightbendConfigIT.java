package smecalculus.rolevod.construction;

import org.springframework.test.context.TestPropertySource;

@TestPropertySource(properties = {"solution.config.mapping.mode=lightbend_config"})
class StorageConfigBeansLightbendConfigIT extends StorageConfigBeansIT {}
