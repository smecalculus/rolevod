package smecalculus.rolevod.storage;

import org.springframework.test.context.ContextConfiguration;
import smecalculus.rolevod.construction.SepulkaDaoBeans;
import smecalculus.rolevod.construction.StoragePropsBeans;

@ContextConfiguration(classes = {StoragePropsBeans.MyBatisPostgres.class, SepulkaDaoBeans.MyBatis.class})
public class SepulkaDaoMyBatisPostgresIT extends SepulkaDaoIT {}
