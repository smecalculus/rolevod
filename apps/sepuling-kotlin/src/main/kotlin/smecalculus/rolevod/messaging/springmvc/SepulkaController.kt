package smecalculus.rolevod.messaging.springmvc

import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController
import smecalculus.rolevod.messaging.SepulkaClient
import smecalculus.rolevod.messaging.SepulkaMessageEm.RegistrationRequest
import smecalculus.rolevod.messaging.SepulkaMessageEm.RegistrationResponse

@RestController
@RequestMapping("sepulkas")
class SepulkaController(
    private val client: SepulkaClient,
) {
    @PostMapping
    fun register(
        @RequestBody requestEdge: RegistrationRequest,
    ): ResponseEntity<RegistrationResponse> {
        val responseEdge = client.register(requestEdge)
        return ResponseEntity.status(HttpStatus.CREATED).body(responseEdge)
    }
}
