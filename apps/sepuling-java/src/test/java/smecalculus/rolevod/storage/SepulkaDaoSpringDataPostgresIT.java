package smecalculus.rolevod.storage;

import org.springframework.test.context.ContextConfiguration;
import smecalculus.rolevod.construction.SepulkaDaoBeans;
import smecalculus.rolevod.construction.StoragePropsBeans;

@ContextConfiguration(classes = {StoragePropsBeans.SpringDataPostgres.class, SepulkaDaoBeans.SpringData.class})
public class SepulkaDaoSpringDataPostgresIT extends SepulkaDaoIT {}
