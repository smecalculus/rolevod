package smecalculus.bezmen.storage;

import org.springframework.test.context.ContextConfiguration;
import smecalculus.bezmen.construction.SepulkaDaoBeans;
import smecalculus.bezmen.construction.StoragePropsBeans;

@ContextConfiguration(classes = {StoragePropsBeans.SpringDataH2.class, SepulkaDaoBeans.SpringData.class})
public class SepulkaDaoSpringDataH2IT extends SepulkaDaoIT {}
