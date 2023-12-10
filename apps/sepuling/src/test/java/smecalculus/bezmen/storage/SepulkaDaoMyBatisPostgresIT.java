package smecalculus.bezmen.storage;

import org.springframework.test.context.ContextConfiguration;
import smecalculus.bezmen.construction.SepulkaDaoBeans;
import smecalculus.bezmen.construction.StoragePropsBeans;

@ContextConfiguration(classes = {StoragePropsBeans.MyBatisPostgres.class, SepulkaDaoBeans.MyBatis.class})
public class SepulkaDaoMyBatisPostgresIT extends SepulkaDaoIT {}
