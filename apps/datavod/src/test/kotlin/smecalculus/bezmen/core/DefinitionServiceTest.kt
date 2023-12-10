package smecalculus.bezmen.core

import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.assertThrows

class DefinitionServiceTest {

    private lateinit var service: DefinitionService

    @BeforeEach
    fun setUp() {
        service = DefinitionService()
    }

    @Test
    fun shouldAssignWithSuccess() {
        // given
        service.register(SuperSystem::class)
        val c1 = service.create("super-system")
        // and
        service.register(SuperCuts::class)
        val c2 = service.create("super-cuts")
        // when
        val assignment = service.assign(c1, "p1", c2)
        // then
        assertThat(assignment.whole).isEqualTo(c1)
        assertThat(assignment.placement).isEqualTo(SuperSystem::subsystem)
        assertThat(assignment.part).isEqualTo(c2)
    }

    @Test
    fun shouldAssignWithException() {
        // given
        service.register(SuperSystem::class)
        val p1 = service.create("super-system")
        // and
        service.register(SubSystem::class)
        val p2 = service.create("sub-system")
        // when
        val exception = assertThrows<IllegalArgumentException> {
            service.assign(p1, "p1", p2)
        }
        // then
        assertThat(exception).hasMessageContaining("must implement")
    }

    @Test
    fun shouldInitiateWithSuccess() {
        // given
        service.register(SuperCuts::class)
        val c1 = service.create("super-cuts")
        // and
        service.register(SubSystem::class)
        val c2 = service.create("sub-system")
        // and
        val assignment = service.assign(c1, "p1", c2)
        // when
        val session = service.initiate(assignment, "bar")
        // then
        assertThat(session.state::class).isEqualTo(TheirChoice::class)
        // when
        service.send(c1, session, Uncut)
        // then
        assertThat(session.state::class).isEqualTo(WeFromThem::class)
        // when
        val response1 = service.receive(c2, session, Uncut::class)
        // then
        assertThat(session.state::class).isEqualTo(OurChoice::class)
        // when
        service.send(c2, session, Cut)
        // then
        assertThat(session.state::class).isEqualTo(TheyFromUs::class)
        // when
        val response2 = service.receive(c1, session, Cut::class)
        // then
        assertThat(session.state::class).isEqualTo(End::class)
    }

    @Test
    fun shouldRestoreWithSuccess() {
        // given
        service.register(SuperCuts::class)
        val c1 = service.create("super-cuts")
        // and
        service.register(SubSystem::class)
        val c2 = service.create("sub-system")
        // and
        val assignment = service.assign(c1, "p1", c2)
        // and
        val expectedSession = service.initiate(assignment, "bar")
        // and
        service.send(c1, expectedSession, Uncut)
        // when
        val actualSession = service.restore(expectedSession.id)
        // then
        assertThat(actualSession)
            .usingRecursiveComparison()
            .isEqualTo(expectedSession)
    }
}

@Role("super-role")
interface SuperAgentRole

@Agent("super-system")
abstract class SuperSystem : SuperAgentRole {
    @Placement("p1")
    abstract val subsystem: AgentOfInterestRole
}

@Role("role-of-interest")
interface AgentOfInterestRole

@Agent("super-cuts")
abstract class SuperCuts : AgentOfInterestRole {
    @Placement("p1")
    abstract val subsystem: Barber
}

@Agent("sub-system")
abstract class SubSystem : Barber

@Alpha("hair")
sealed class Hair

@State("uncut")
data object Uncut : Hair()

@State("cut")
data object Cut : Hair()

@Role("barber")
interface Barber {
    companion object : Barber

    @Port("bar")
    fun bar(): Interplay {
        return interplay("haircut") {
            their<Hair> {
                choice(Cut::class) {
                    end()
                }
                choice(Uncut::class) {
                    our<Hair> {
                        choice(Cut::class) {
                            end()
                        }
                    }
                }
            }
        }
    }

    @Port("ber")
    fun ber(): Interplay {
        return interplay("hairstyle") {}
    }
}
