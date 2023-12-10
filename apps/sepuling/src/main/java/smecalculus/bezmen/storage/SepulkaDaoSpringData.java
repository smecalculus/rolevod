package smecalculus.bezmen.storage;

import java.util.Optional;
import java.util.UUID;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import smecalculus.bezmen.core.StateDm;
import smecalculus.bezmen.storage.springdata.SepulkaRepository;

@RequiredArgsConstructor
public class SepulkaDaoSpringData implements SepulkaDao {

    @NonNull
    private SepulkaStateMapper mapper;

    @NonNull
    private SepulkaRepository repository;

    @Override
    public StateDm.AggregateState add(@NonNull StateDm.AggregateState state) {
        var stateEdge = repository.save(mapper.toEdge(state));
        return mapper.toDomain(stateEdge);
    }

    @Override
    public Optional<StateDm.ExistenceState> getBy(@NonNull String externalId) {
        return repository
                .findByExternalId(externalId, StateEm.ExistenceState.class)
                .map(mapper::toDomain);
    }

    @Override
    public Optional<StateDm.PreviewState> getBy(@NonNull UUID internalId) {
        return repository
                .findByInternalId(internalId.toString(), StateEm.PreviewState.class)
                .map(mapper::toDomain);
    }

    @Override
    public void updateBy(StateDm.TouchState state, UUID internalId) {
        var stateEdge = mapper.toEdge(state);
        var matchedCount = repository.updateBy(stateEdge, internalId.toString());
        if (matchedCount == 0) {
            throw new ContentionException();
        }
    }
}
