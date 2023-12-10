package smecalculus.bezmen.core

import smecalculus.bezmen.core.SepulkaMessageDm.PreviewRequest
import smecalculus.bezmen.core.SepulkaMessageDm.PreviewResponse
import smecalculus.bezmen.core.SepulkaMessageDm.RegistrationRequest
import smecalculus.bezmen.core.SepulkaMessageDm.RegistrationResponse
import smecalculus.bezmen.storage.SepulkaDao
import java.time.LocalDateTime
import java.util.UUID

class SepulkaServiceImpl(
    private val mapper: SepulkaMapper,
    private val dao: SepulkaDao,
) : SepulkaService {
    override fun register(request: RegistrationRequest): RegistrationResponse {
        val now = LocalDateTime.now()
        val sepulkaCreated =
            mapper.toState(request)
                .internalId(UUID.randomUUID())
                .revision(0)
                .createdAt(now)
                .updatedAt(now)
                .build()
        val sepulkaSaved = dao.add(sepulkaCreated)
        return mapper.toMessage(sepulkaSaved).build()
    }

    override fun view(request: PreviewRequest): List<PreviewResponse> {
        return listOf()
    }
}
