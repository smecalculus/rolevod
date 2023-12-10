package smecalculus.rolevod.messaging

import smecalculus.rolevod.core.SepulkaService
import smecalculus.rolevod.messaging.SepulkaMessageEm.RegistrationRequest
import smecalculus.rolevod.messaging.SepulkaMessageEm.RegistrationResponse
import smecalculus.rolevod.validation.EdgeValidator

class SepulkaClientImpl(
    private val validator: EdgeValidator,
    private val mapper: SepulkaMessageMapper,
    private val service: SepulkaService,
) : SepulkaClient {
    override fun register(requestEdge: RegistrationRequest): RegistrationResponse {
        validator.validate(requestEdge)
        val request = mapper.toDomain(requestEdge)
        val response = service.register(request)
        return mapper.toEdge(response)
    }
}
