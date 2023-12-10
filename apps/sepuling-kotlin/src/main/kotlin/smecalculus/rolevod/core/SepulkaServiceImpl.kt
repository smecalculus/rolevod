package smecalculus.rolevod.core

import smecalculus.rolevod.core.SepulkaMessageDm.PreviewRequest
import smecalculus.rolevod.core.SepulkaMessageDm.PreviewResponse
import smecalculus.rolevod.core.SepulkaMessageDm.RegistrationRequest
import smecalculus.rolevod.core.SepulkaMessageDm.RegistrationResponse
import smecalculus.rolevod.storage.SepulkaDao
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
