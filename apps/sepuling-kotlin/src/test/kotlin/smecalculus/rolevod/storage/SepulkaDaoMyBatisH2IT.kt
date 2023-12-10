package smecalculus.rolevod.storage

import org.springframework.test.context.ContextConfiguration
import smecalculus.rolevod.construction.SepulkaDaoBeans
import smecalculus.rolevod.construction.StoragePropsBeans

@ContextConfiguration(classes = [StoragePropsBeans.MyBatisH2::class, SepulkaDaoBeans.MyBatis::class])
class SepulkaDaoMyBatisH2IT : SepulkaDaoIT()
