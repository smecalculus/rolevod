package smecalculus.rolevod.storage;

import org.springframework.test.context.ContextConfiguration;
import smecalculus.rolevod.construction.SepulkaDaoBeans;
import smecalculus.rolevod.construction.StoragePropsBeans;

@ContextConfiguration(classes = {StoragePropsBeans.SpringDataH2.class, SepulkaDaoBeans.SpringData.class})
public class SepulkaDaoSpringDataH2IT extends SepulkaDaoIT {}
