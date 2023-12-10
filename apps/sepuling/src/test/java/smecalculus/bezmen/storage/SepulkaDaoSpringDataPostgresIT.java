package smecalculus.bezmen.storage;

import org.springframework.test.context.ContextConfiguration;
import smecalculus.bezmen.construction.SepulkaDaoBeans;
import smecalculus.bezmen.construction.StoragePropsBeans;

@ContextConfiguration(classes = {StoragePropsBeans.SpringDataPostgres.class, SepulkaDaoBeans.SpringData.class})
public class SepulkaDaoSpringDataPostgresIT extends SepulkaDaoIT {}
