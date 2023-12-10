package smecalculus.bezmen.storage

import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.extension.ExtendWith
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.test.annotation.DirtiesContext
import org.springframework.test.context.ContextConfiguration
import org.springframework.test.context.junit.jupiter.SpringExtension
import smecalculus.bezmen.construction.SepulkaDaoBeans
import smecalculus.bezmen.core.SepulkaStateDmEg

@DirtiesContext
@ExtendWith(SpringExtension::class)
@ContextConfiguration(classes = [SepulkaDaoBeans.Anyone::class])
abstract class SepulkaDaoIT {
    @Autowired
    private lateinit var sepulkaDao: SepulkaDao

    @Test
    fun shouldAddOneSepulka() {
        // given
        val expected1 = SepulkaStateDmEg.aggregateRoot().build()
        // and
        val expected2 = SepulkaStateDmEg.existence().internalId(expected1.internalId).build()
        // when
        val actualSaved = sepulkaDao.add(expected1)
        // and
        val actualSelected = sepulkaDao.getBy(expected1.externalId)
        // then
        assertThat(actualSaved).usingRecursiveComparison().isEqualTo(expected1)
        // and
        assertThat(actualSelected).usingRecursiveComparison().isEqualTo(expected2)
    }

    @Test
    fun shouldViewOneSepulka() {
        // given
        val aggregate = SepulkaStateDmEg.aggregateRoot().build()
        // and
        sepulkaDao.add(aggregate)
        // and
        val expected = SepulkaStateDmEg.preview(aggregate).build()
        // when
        val actual = sepulkaDao.getBy(aggregate.internalId)
        // then
        assertThat(actual).isEqualTo(expected)
    }

    @Test
    fun shouldUpdateOneSepulka() {
        // given
        val aggregate = SepulkaStateDmEg.aggregateRoot().build()
        // and
        sepulkaDao.add(aggregate)
        // and
        val updatedAt = aggregate.updatedAt.plusSeconds(1)
        // and
        val touch = SepulkaStateDmEg.touch(aggregate).updatedAt(updatedAt).build()
        // when
        sepulkaDao.updateBy(aggregate.internalId, touch)
        // then
        // no exception
    }
}
